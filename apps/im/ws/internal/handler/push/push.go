/**
 * @author: Yanko/xiaoxiaoyang-sheep
 * @doc:
 **/

package push

import (
	"easy-chat/apps/im/ws/internal/svc"
	"easy-chat/apps/im/ws/websocket"
	"easy-chat/apps/im/ws/ws"
	"github.com/mitchellh/mapstructure"
)

func Push(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		var data *ws.Push
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err))
			return
		}

		// 发送的目标
		rconn := srv.GetConn(data.RecvId)
		if rconn == nil {
			// TODO: 目标离线
			return
		}

		srv.Info("recv push msg %v", data)

		srv.Send(websocket.NewMessage(data.SendId, &ws.Chat{
			ConversationId: data.ConversationId,
			ChatType:       data.ChatType,
			Msg: ws.Msg{
				MType:   data.MType,
				Content: data.Content,
			},
			SendTime: data.SendTime,
		}), rconn)
	}
}
