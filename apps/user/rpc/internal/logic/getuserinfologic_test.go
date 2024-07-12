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

func TestGetUserInfoLogic_GetUserInfo(t *testing.T) {

	type args struct {
		in *user.GetUserInfoReq
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"getUserInfo", args{in: &user.GetUserInfoReq{
				Id: "0x0000001000000001",
			}}, true, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewGetUserInfoLogic(context.Background(), svcCtx)
			got, err := l.GetUserInfo(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want {
				t.Log(tt.name, got)
			}
		})
	}
}
