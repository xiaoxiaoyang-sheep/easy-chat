/**
 * @author: Yanko/xiaoxiaoyang-sheep
 * @doc:
 **/

package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/url"
)

type Client interface {
	Close() error

	Send(v interface{}) error
	Rend(v interface{}) error
}

type client struct {
	*websocket.Conn
	host string
	opt  dailOption
}

func NewClient(host string, opts ...DailOptions) *client {
	opt := NewDailOptions(opts...)

	c := client{
		Conn: nil,
		host: host,
		opt:  opt,
	}

	conn, err := c.dail()
	if err != nil {
		panic(err)
	}

	c.Conn = conn
	return &c
}

func (c *client) dail() (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws", Host: c.host, Path: c.opt.pattern}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), c.opt.header)
	return conn, err
}

func (c *client) Send(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	err = c.Conn.WriteMessage(websocket.TextMessage, data)
	if err == nil {
		return nil
	}

	// todo: 在增加一个重连发送
	conn, err := c.dail()
	if err != nil {
		return err
	}
	c.Conn = conn
	return c.Conn.WriteMessage(websocket.TextMessage, data)
}

func (c *client) Rend(v interface{}) error {
	_, msg, err := c.Conn.ReadMessage()
	if err != nil {
		return err
	}

	return json.Unmarshal(msg, v)
}
