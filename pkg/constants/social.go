/**
 * @author: Yanko/xiaoxiaoyang-sheep
 * @doc:
 **/

package constants

// 处理结构 1. 未处理  2. 处理  3. 拒绝
type HandlerResult int

const (
	NoHandlerResult HandlerResult = iota + 1
	PassHandlerResult
	RefuseHandlerResult
	CancelHandlerResult
)
