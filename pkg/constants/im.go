/**
 * @author: Yanko/xiaoxiaoyang-sheep
 * @doc:
 **/

package constants

type MType int

const (
	TextMtype MType = iota
)

type ChatType int

const (
	GroupChatType ChatType = iota + 1
	SingleChatType
)
