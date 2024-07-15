/**
 * @author: Yanko/xiaoxiaoyang-sheep
 * @doc:
 **/

package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

type Conn struct {
	idleMu sync.Mutex
	*websocket.Conn
	Uid               string
	s                 *Server
	idle              time.Time
	maxConnectionIdle time.Duration

	messageMu      sync.Mutex
	readMessage    []*Message
	readMessageSeq map[string]*Message
	message        chan *Message

	done chan struct{}
}

func NewConn(s *Server, w http.ResponseWriter, r *http.Request) *Conn {
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Errorf("upgrade err &v", err)
		return nil
	}

	conn := &Conn{
		Conn:              c,
		s:                 s,
		idle:              time.Now(),
		maxConnectionIdle: s.opt.maxConnectionIdle,

		message:        make(chan *Message, 1),
		readMessage:    make([]*Message, 0, 2),
		readMessageSeq: make(map[string]*Message, 2),

		done: make(chan struct{}),
	}

	go conn.keepalive()
	return conn
}

func (c *Conn) ReadMessage() (messageType int, p []byte, err error) {
	messageType, p, err = c.Conn.ReadMessage()
	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	c.idle = time.Now()
	return
}

func (c *Conn) appendMsgMq(msg *Message) {
	c.messageMu.Lock()
	defer c.messageMu.Unlock()

	// 读队列中
	if m, ok := c.readMessageSeq[msg.Id]; ok {
		// 已经有消息的记录，该消息已经有ack确认
		if len(c.readMessage) == 0 {
			// 队列中没有该消息
			return
		}

		if m.AckSeq >= msg.AckSeq {
			// 没有进行ack确认，或者重复
			return
		}

		c.readMessageSeq[msg.Id] = msg
		return
	}

	// 还没有进行ack确认，避免客户端重复发送多余的ack消息
	if msg.FrameType == FrameAck {
		return
	}

	c.readMessage = append(c.readMessage, msg)
	c.readMessageSeq[msg.Id] = msg
}

func (c *Conn) Close() error {
	select {
	case <-c.done:
	default:
		close(c.done)
	}

	return c.Conn.Close()
}

func (c *Conn) WriteMessage(messageType int, data []byte) error {
	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	err := c.Conn.WriteMessage(messageType, data)
	c.idle = time.Now()
	return err
}

func (c *Conn) keepalive() {
	idleTimer := time.NewTimer(c.maxConnectionIdle)
	defer func() {
		idleTimer.Stop()
	}()

	for {
		select {
		case <-idleTimer.C:
			c.idleMu.Lock()
			idle := c.idle
			if idle.IsZero() {
				c.idleMu.Unlock()
				idleTimer.Reset(c.maxConnectionIdle)
				continue
			}
			val := c.maxConnectionIdle - time.Since(idle)
			c.idleMu.Unlock()
			if val <= 0 {
				c.s.Close(c)
				return
			}
			idleTimer.Reset(val)
		case <-c.done:
			return
		}
	}
}
