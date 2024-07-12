package logic

import (
	"context"
	"easy-chat/apps/user/models"
	"easy-chat/pkg/xerr"
	"github.com/jinzhu/copier"

	"easy-chat/apps/user/rpc/internal/svc"
	"easy-chat/apps/user/rpc/user"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLogic {
	return &FindUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindUserLogic) FindUser(in *user.FindUserReq) (*user.FindUserResp, error) {
	// todo: add your logic here and delete this line

	var (
		userEntitys []*models.Users
		err         error
	)

	if in.Phone != "" {
		var userEntity *models.Users
		userEntity, err = l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
		if err == nil {
			userEntitys = append(userEntitys, userEntity)
		}
	} else if in.Name != "" {
		userEntitys, err = l.svcCtx.UsersModel.ListByName(l.ctx, in.Name)
	} else if len(in.Ids) > 0 {
		userEntitys, err = l.svcCtx.UsersModel.ListByIds(l.ctx, in.Ids)
	}

	if err != nil && !errors.Is(err, models.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find by phone or name or ids err %v, req %v",
			err, *in)
	}

	var resp []*user.UserEntity
	copier.Copy(&resp, &userEntitys)

	return &user.FindUserResp{
		User: resp,
	}, nil
}
