/**
 * @author: Yanko/xiaoxiaoyang-sheep
 * @doc:
 **/

package conversation

import (
	"context"
	"easy-chat/apps/im/ws/internal/logic"
	"easy-chat/apps/im/ws/internal/svc"
	"easy-chat/apps/im/ws/websocket"
	"easy-chat/apps/im/ws/ws"
	"easy-chat/pkg/constants"
	"github.com/mitchellh/mapstructure"
	"time"
)

func Chat(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		// todo: 私聊

		var data ws.Chat
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}

		switch data.ChatType {
		case constants.SingleChatType:
			err := logic.NewConversation(context.Background(), srv, svc).SingleChat(&data,
				conn.Uid)
			if err != nil {
				srv.Send(websocket.NewErrMessage(err), conn)
				return
			}

			srv.SendByUserId(websocket.NewMessage(conn.Uid, ws.Chat{
				ConversationId: data.ConversationId,
				ChatType:       data.ChatType,
				SendId:         conn.Uid,
				RecvId:         data.RecvId,
				Msg:            data.Msg,
				SendTime:       time.Now().UnixMilli(),
			}), data.RecvId)
		}
	}
}
