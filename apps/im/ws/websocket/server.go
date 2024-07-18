/**
 * @author: Yanko/xiaoxiaoyang-sheep
 * @doc:
 **/

package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"net/http"
	"sync"
	"time"
)

type AckType int

const (
	NoAck AckType = iota
	OnlyAck
	RigorAck
)

func (t AckType) ToString() string {
	switch t {
	case OnlyAck:
		return "OnlyAck"
	case RigorAck:
		return "RigorAck"
	}
	return "NoAck"
}

type Server struct {
	sync.RWMutex
	opt            *serverOption
	authentication Authentication
	patten         string
	routes         map[string]HandlerFunc
	addr           string
	connToUser     map[*Conn]string
	userToConn     map[string]*Conn
	upgrader       websocket.Upgrader
	*threading.TaskRunner
	logx.Logger
}

func NerServer(addr string, opts ...ServerOptions) *Server {
	opt := newServerOptions(opts...)

	return &Server{
		routes:         make(map[string]HandlerFunc),
		addr:           addr,
		opt:            &opt,
		patten:         opt.patten,
		authentication: opt.Authentication,

		connToUser: make(map[*Conn]string),
		userToConn: make(map[string]*Conn),

		upgrader:   websocket.Upgrader{},
		TaskRunner: threading.NewTaskRunner(opt.concurrent),
		Logger:     logx.WithContext(context.Background()),
	}
}

func (s *Server) ServerWs(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			s.Errorf("server handler ws recover err %v", r)
		}
	}()

	conn := NewConn(s, w, r)
	if conn == nil {
		return
	}

	// 鉴权
	if !s.authentication.Auth(w, r) {
		s.Send(&Message{FrameType: FrameData, Data: fmt.Sprintf("不具备访问权限")}, conn)
		conn.Close()
		return
	}

	// 记录连接
	s.addConn(conn, r)

	// 处理连接
	go s.handlerConn(conn)
}

func (s *Server) addConn(conn *Conn, r *http.Request) {
	uid := s.authentication.UserId(r)

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	// 验证用户是否之前登录过
	if c := s.userToConn[uid]; c != nil {
		// 关闭之前的连接
		c.Close()
	}

	s.connToUser[conn] = uid
	s.userToConn[uid] = conn
}

// 根据连接对象执行任务处理
func (s *Server) handlerConn(conn *Conn) {

	uids := s.GetUsers(conn)
	conn.Uid = uids[0]

	// 处理任务
	go s.handlerWrite(conn)

	if s.isAck(nil) {
		go s.ReadAck(conn)
	}

	for {
		// 获取请求消息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Errorf("websocket conn read message err %v", err)
			s.Close(conn)
			return
		}

		// 解析消息
		var message Message
		if err = json.Unmarshal(msg, &message); err != nil {
			s.Errorf("json unmarshal err %v, msg %v", err, string(msg))
			s.Close(conn)
			return
		}

		// 是否需要ack验证
		if s.isAck(&message) {
			s.Info("conn message read ack msg %v", message)
			conn.AppendMsgMq((&message))
		} else {
			conn.message <- &message
		}
	}
}

func (s *Server) isAck(message *Message) bool {
	if message == nil {
		return s.opt.ack != NoAck
	}
	return s.opt.ack != NoAck && message.FrameType != FrameNoAck
}

// 读取消息的ack
func (s *Server) ReadAck(conn *Conn) {

	send := func(msg *Message, conn *Conn) error {
		err := s.Send(msg, conn)
		if err == nil {
			return nil
		}

		s.Errorf("message ack OnlyAck send err %v message %v", err, msg)
		conn.readMessage[0].errCount++
		conn.messageMu.Unlock()

		tempDelay := time.Duration(200*conn.readMessage[0].errCount) * time.Microsecond
		if max := 1 * time.Second; tempDelay > max {
			tempDelay = max
		}

		time.Sleep(tempDelay)
		return err
	}

	for {
		select {
		case <-conn.done:
			s.Info("close message ack uid %v", conn.Uid)
			return
		default:
		}

		// 从队列中读取新的消息
		conn.messageMu.Lock()
		if len(conn.readMessage) == 0 {
			conn.messageMu.Unlock()
			// 增加睡眠
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// 读取第一条
		message := conn.readMessage[0]

		// 判断ack的方式
		switch s.opt.ack {
		case OnlyAck:
			// 给客户端发送消息
			if message.FormId != conn.Uid {
				if message.AckSeq == -1 {
					conn.readMessage[0].AckSeq++
					conn.readMessage[0].ackTime = time.Now()
					if err := send(message, conn); err != nil {
						continue
					}
				}

				msgSeq := conn.readMessageSeq[message.Id]
				if msgSeq.AckSeq == 1 {
					conn.readMessage = conn.readMessage[1:]
					delete(conn.readMessageSeq, message.Id)
					conn.messageMu.Unlock()
					s.Infof("message send to client success mid %v", message.Id)
					continue
				}

				// 客户端没有确认，考虑是否超过了ack的确认时间
				val := s.opt.ackTimeout - time.Since(message.ackTime)
				if !message.ackTime.IsZero() && val <= 0 {
					//  超过 结束确认
					// TODO: 超时了，可以选择断开与客户端的连接,但实际具体细节处理仍然还需自己结合业务完善，此处选择放弃该消息
					s.Errorf("message ack RigorAck fail mid %v, time %v because timeout", message.Id, message.ackTime)
					delete(conn.readMessageSeq, message.Id)
					conn.readMessage = conn.readMessage[1:]
					conn.messageMu.Unlock()
					continue
				}

				// 未超过 重新发送
				conn.messageMu.Unlock()
				// 避免第一次确认后重复发送
				if val > 0 && val <= s.opt.ackTimeout-3000*time.Millisecond {
					if err := send(message, conn); err != nil {
						continue
					}
				}
				// 睡眠一定的时间
				time.Sleep(3000 * time.Millisecond)

			} else {
				// 接收客户端消息
				// 直接给客户端回复
				if err := send(&Message{
					FrameType: FrameAck,
					Id:        message.Id,
					AckSeq:    message.AckSeq + 1,
				}, conn); err != nil {
					continue
				}
				// 进行业务处理
				// 把消息从队列中移除
				conn.readMessage = conn.readMessage[1:]
				conn.messageMu.Unlock()
				conn.message <- message
				s.Infof("message ack OnlyAck send success mid %v", message.Id)
			}
		case RigorAck:

			// 给客户端发送消息
			if message.FormId != conn.Uid {
				if message.AckSeq == -1 {
					conn.readMessage[0].AckSeq++
					conn.readMessage[0].ackTime = time.Now()
					if err := send(message, conn); err != nil {
						continue
					}
				}

				msgSeq := conn.readMessageSeq[message.Id]
				if msgSeq.AckSeq == 1 {
					if err := send(&Message{
						FrameType: FrameAck,
						AckSeq:    2,
						Id:        message.Id,
					}, conn); err != nil {
						continue
					}
					conn.readMessage = conn.readMessage[1:]
					delete(conn.readMessageSeq, message.Id)
					conn.messageMu.Unlock()
					s.Infof("message send to client ack success", message.Id)
					continue
				}

				// 客户端没有确认，考虑是否超过了ack的确认时间
				val := s.opt.ackTimeout - time.Since(message.ackTime)
				if !message.ackTime.IsZero() && val <= 0 {
					//  超过 结束确认
					// TODO: 超时了，可以选择断开与客户端的连接,但实际具体细节处理仍然还需自己结合业务完善，此处选择放弃该消息
					s.Errorf("message ack RigorAck fail mid %v, time %v because timeout", message.Id, message.ackTime)
					delete(conn.readMessageSeq, message.Id)
					conn.readMessage = conn.readMessage[1:]
					conn.messageMu.Unlock()
					continue
				}

				//  未超过 重新发送
				conn.messageMu.Unlock()
				// 避免第一次确认后重复发送
				if val > 0 && val <= s.opt.ackTimeout-3000*time.Millisecond {
					if err := send(message, conn); err != nil {
						continue
					}
				}
				// 睡眠一定的时间
				time.Sleep(3000 * time.Millisecond)
			} else {
				// 接收客户端消息

				// 先回
				if message.AckSeq == 0 {
					// 还未确认
					conn.readMessage[0].AckSeq++
					conn.readMessage[0].ackTime = time.Now()
					if err := send(&Message{
						FrameType: FrameAck,
						AckSeq:    message.AckSeq,
						Id:        message.Id,
					}, conn); err != nil {
						continue
					}

					s.Info("message ack RigorAck send mid %v, seq %v, time %v", message.Id,
						message.AckSeq, message.ackTime)
					conn.messageMu.Unlock()
					continue
				}

				// 再验证
				// 1. 客户端返回结果，再一次确认
				// 得到客户端的序号
				msgSeq := conn.readMessageSeq[message.Id]

				if msgSeq.AckSeq > message.AckSeq {
					// 确认
					conn.readMessage = conn.readMessage[1:]
					conn.messageMu.Unlock()
					conn.message <- message
					s.Infof("message ack RigorAck sucess mid %v", message.Id)
					continue
				}

				// 2. 客户端没有确认，考虑是否超过了ack的确认时间
				val := s.opt.ackTimeout - time.Since(message.ackTime)
				if !message.ackTime.IsZero() && val <= 0 {
					//  2.2 超过 结束确认
					// TODO: 超时了，可以选择断开与客户端的连接,但实际具体细节处理仍然还需自己结合业务完善，此处选择放弃该消息
					s.Errorf("message ack RigorAck fail mid %v, time %v because timeout", message.Id, message.ackTime)
					delete(conn.readMessageSeq, message.Id)
					conn.readMessage = conn.readMessage[1:]
					conn.messageMu.Unlock()
					continue
				}

				//  2.1 未超过 重新发送
				conn.messageMu.Unlock()
				// 避免第一次确认后重复发送
				if val > 0 && val <= s.opt.ackTimeout-3000*time.Millisecond {
					if err := send(&Message{
						FrameType: FrameAck,
						AckSeq:    message.AckSeq,
						Id:        message.Id,
					}, conn); err != nil {
						continue
					}
				}
				// 睡眠一定的时间
				time.Sleep(3000 * time.Millisecond)
			}

		}

	}
}

// 任务的处理
func (s *Server) handlerWrite(conn *Conn) {
	for {
		select {
		case <-conn.done:
			// 通道关闭
			return
		case message := <-conn.message:
			// 依据消息进行处理
			switch message.FrameType {
			case FramePing:
				s.Send(&Message{FrameType: FramePing, FormId: conn.Uid}, conn)
			case FrameData:
				fallthrough
			case FrameNoAck:
				// 根据请求的method分发路由并执行
				if handler, ok := s.routes[message.Method]; ok {
					handler(s, conn, message)
				} else {
					s.Send(&Message{FrameType: FrameData, Data: fmt.Sprintf("不存在执行的方法 %v, 请检查",
						message.Method)}, conn)
				}
			}

			if s.isAck(message) {
				conn.messageMu.Lock()
				delete(conn.readMessageSeq, message.Id)
				conn.messageMu.Unlock()
			}
		}
	}
}

func (s *Server) GetConn(uid string) *Conn {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	return s.userToConn[uid]
}

func (s *Server) GetConns(uids ...string) []*Conn {
	if len(uids) == 0 {
		return nil
	}
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	res := make([]*Conn, 0, len(uids))
	for _, uid := range uids {
		res = append(res, s.userToConn[uid])
	}

	return res
}

func (s *Server) GetUsers(conns ...*Conn) []string {

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	var res []string
	if len(conns) == 0 {
		// 获取全部
		res = make([]string, 0, len(s.connToUser))
		for _, uid := range s.connToUser {
			res = append(res, uid)
		}
	} else {
		// 获取部分
		res = make([]string, 0, len(conns))
		for _, conn := range conns {
			res = append(res, s.connToUser[conn])
		}
	}

	return res
}

func (s *Server) Close(conn *Conn) {

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	uid := s.connToUser[conn]
	if uid == "" {
		// 已经关闭
		return
	}
	delete(s.connToUser, conn)
	delete(s.userToConn, uid)

	conn.Close()
}

func (s *Server) Send(msg interface{}, conns ...*Conn) error {
	if len(conns) == 0 {
		return nil
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	for _, conn := range conns {
		if err = conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) SendByUserId(msg interface{}, userIds ...string) error {
	if len(userIds) == 0 {
		return nil
	}

	return s.Send(msg, s.GetConns(userIds...)...)
}

func (s *Server) AddRoutes(rs []Route) {
	for _, r := range rs {
		s.routes[r.Method] = r.Handler
	}
}

func (s *Server) Start() {
	http.HandleFunc(s.patten, s.ServerWs)
	s.Info(http.ListenAndServe(s.addr, nil))
}

func (s *Server) Stop() {
	fmt.Println("停止服务")
}

func (s *Server) GetOptAct() AckType {
	return s.opt.ack
}
