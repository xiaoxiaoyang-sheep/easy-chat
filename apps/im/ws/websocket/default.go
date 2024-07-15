/**
 * @author: Yanko/xiaoxiaoyang-sheep
 * @doc:
 **/

package websocket

import (
	"math"
	"time"
)

const (
	defaultPattern           = "/ws"
	defaultMaxConnectionIdle = time.Duration(math.MaxInt64)
	defaultAckTimeout        = 30 * time.Second
	defaultAck               = NoAck
)
