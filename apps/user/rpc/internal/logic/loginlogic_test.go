/**
 * @author: Yanko/xiaoxiaoyang-sheep
 * @doc:
 **/

package logic

import (
	"context"
	"easy-chat/apps/user/rpc/user"
	"testing"
)

func TestLoginLogic_Login(t *testing.T) {

	type args struct {
		in *user.LoginReq
	}
	tests := []struct {
		name      string
		args      args
		wantPrint bool
		wantErr   bool
	}{
		{"login", args{in: &user.LoginReq{
			Phone:    "18758004742",
			Password: "1234567",
		}}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLoginLogic(context.Background(), svcCtx)
			got, err := l.Login(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantPrint {
				t.Log(tt.name, got)
			}
		})
	}
}
