/**
 * @author: Yanko/xiaoxiaoyang-sheep
 * @doc:
 **/

package mqclient

import (
	"easy-chat/apps/task/mq/mq"
	"encoding/json"
	"github.com/zeromicro/go-queue/kq"
)

type MsgChatTransferClient interface {
	Push(msg *mq.TaskChatTransfer) error
}

type msgChatTransferClient struct {
	pusher *kq.Pusher
}

func NewMsgChatTransferClient(addr []string, topic string, opts ...kq.PushOption) MsgChatTransferClient {
	return &msgChatTransferClient{
		pusher: kq.NewPusher(addr, topic),
	}
}

func (c *msgChatTransferClient) Push(msg *mq.TaskChatTransfer) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return c.pusher.Push(string(body))
}
