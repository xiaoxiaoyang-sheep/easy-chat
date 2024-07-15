/**
 * @author: Yanko/xiaoxiaoyang-sheep
 * @doc:
 **/

package ws

import (
	"easy-chat/pkg/constants"
)

type (
	Msg struct {
		constants.MType `mapstructure:"mType"`
		Content         string `mapstructure:"content"`
	}

	Chat struct {
		ConversationId     string `mapstructure:"conversationId"`
		constants.ChatType `mapstructure:"chatType"`
		SendId             string `mapstructure:"sendId" json:"sendId,omitempty"`
		RecvId             string `mapstructure:"recvId" json:"recvId,omitempty"`
		Msg                `mapstructure:"msg"`
		SendTime           int64 `mapstructure:"sendTime"`
	}

	Push struct {
		ConversationId     string `mapstructure:"conversationId"`
		constants.ChatType `json:"chatType"`
		SendId             string `mapstructure:"sendId"`
		RecvId             string `mapstructure:"recvId"`
		SendTime           int64  `mapstructure:"sendTime"`

		constants.MType `mapstructure:"mType"`
		Content         string `mapstructure:"content"`
	}
)
