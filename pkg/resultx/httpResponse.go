/**
 * @author: Yanko/xiaoxiaoyang-sheep
 * @doc:
 **/

package resultx

import (
	"context"
	"easy-chat/pkg/xerr"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	zerr "github.com/zeromicro/x/errors"
	"google.golang.org/grpc/status"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(data interface{}) *Response {
	return &Response{
		Code: 200,
		Msg:  "请求成功",
		Data: data,
	}
}

func Fail(code int, err string) *Response {
	return &Response{
		Code: code,
		Msg:  err,
		Data: nil,
	}
}

func OkHandler(_ context.Context, v interface{}) any {
	return Success(v)
}

func ErrHandler(name string) func(ctx context.Context, err error) (int, any) {
	return func(ctx context.Context, err error) (int, any) {
		errcode := xerr.SERVER_COMMON_ERROR
		errmsg := xerr.ErrMsg(errcode)

		causeErr := errors.Cause(err)
		var e *zerr.CodeMsg
		if errors.As(causeErr, &e) {
			errcode = e.Code
			errmsg = e.Msg
		} else {
			if gsstatus, ok := status.FromError(causeErr); ok {
				grpcCode := int(gsstatus.Code())
				errcode = grpcCode
				errmsg = gsstatus.Message()
			}
		}

		// 日志记录
		logx.WithContext(ctx).Errorf("【%s】err %v", name, err)

		return http.StatusBadRequest, Fail(errcode, errmsg)
	}
}
