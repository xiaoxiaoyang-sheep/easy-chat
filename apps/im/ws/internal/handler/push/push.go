/**
 * @author: Yanko/xiaoxiaoyang-sheep
 * @doc:
 **/

package push

import (
	"easy-chat/apps/im/ws/internal/svc"
	"easy-chat/apps/im/ws/websocket"
	"easy-chat/apps/im/ws/ws"
	"easy-chat/pkg/constants"
	"github.com/mitchellh/mapstructure"
	"strconv"
	"time"
)

func Push(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		var data *ws.Push
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err))
			return
		}

		// 发送的目标
		switch data.ChatType {
		case constants.SingleChatType:
			single(srv, data, data.RecvId)
		case constants.GroupChatType:
			group(srv, data)
		}

	}
}

func single(srv *websocket.Server, data *ws.Push, recvId string) error {

	rconn := srv.GetConn(recvId)
	if rconn == nil {
		// TODO: 目标离线
		return nil
	}

	srv.Info("recv push msg %v", data)

	message := websocket.NewMessage(data.SendId, &ws.Chat{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		Msg: ws.Msg{
			MType:   data.MType,
			Content: data.Content,
		},
		SendTime: data.SendTime,
	})

	if srv.GetOptAct() == websocket.NoAck {
		return srv.Send(message, rconn)
	} else {
		message.AckSeq = -1
		message.Id = strconv.FormatInt(time.Now().UnixMilli(), 10)
		rconn.AppendMsgMq(message)
	}

	return nil
}

func group(srv *websocket.Server, data *ws.Push) error {
	for _, id := range data.RecvIds {
		func(id string) {
			srv.Schedule(func() {
				single(srv, data, id)
			})
		}(id)
	}
	return nil
}
