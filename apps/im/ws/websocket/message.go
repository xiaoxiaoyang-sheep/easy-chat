/**
 * @author: Yanko/xiaoxiaoyang-sheep
 * @doc:
 **/

package websocket

import "time"

type FrameType uint8

const (
	FrameData  FrameType = 0x0
	FramePing  FrameType = 0x1
	FrameAck   FrameType = 0x2
	FrameNoAck FrameType = 0x3
	FrameErr   FrameType = 0x9
)

type Message struct {
	FrameType FrameType   `json:"frameType"`
	Id        string      `json:"id,omitempty"`
	AckSeq    int         `json:"ackSeq,omitempty"`
	ackTime   time.Time   `json:"ackTime"`
	errCount  int         `json:"errCount"`
	Method    string      `json:"method,omitempty"`
	FormId    string      `json:"formId,omitempty"`
	Data      interface{} `json:"data"`
}

func NewMessage(formId string, data interface{}) *Message {
	return &Message{
		FrameType: FrameData,
		FormId:    formId,
		Data:      data,
	}
}

func NewErrMessage(err error) *Message {
	return &Message{
		FrameType: FrameErr,
		Data:      err.Error(),
	}
}
