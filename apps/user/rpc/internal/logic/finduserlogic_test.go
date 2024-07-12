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

func TestFindUserLogic_FindUser(t *testing.T) {
	type args struct {
		in *user.FindUserReq
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"phone", args{in: &user.FindUserReq{
				Phone: "18758004743",
			}}, true, false,
		},

		{"name", args{in: &user.FindUserReq{
			Name: "s",
		}}, true, false},

		{
			"ids", args{in: &user.FindUserReq{
				Ids: []string{"0x0000001000000001", "0x0000002000000001", "0x0000003000000001"},
			}}, true, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewFindUserLogic(context.Background(), svcCtx)
			got, err := l.FindUser(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want {
				t.Log(tt.name, got)
			}
		})
	}
}
