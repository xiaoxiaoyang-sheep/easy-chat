package logic

import (
	"context"
	"database/sql"
	"easy-chat/apps/user/models"
	"easy-chat/pkg/ctxdata"
	"easy-chat/pkg/encrypt"
	"easy-chat/pkg/wuid"
	"easy-chat/pkg/xerr"
	"github.com/pkg/errors"
	"time"

	"easy-chat/apps/user/rpc/internal/svc"
	"easy-chat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrPhoneIsRegister = xerr.New(xerr.SERVER_COMMON_ERROR, "手机号已经注册过")
	ErrEncryptionPwd   = xerr.New(xerr.SERVER_COMMON_ERROR, "密码加密出错")
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	// todo: add your logic here and delete this line

	// 验证用户是否注册，根据手机号验证
	userEntity, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find user by phone err %v, req %v", err, in.Phone)
	}

	if userEntity != nil {
		return nil, errors.WithStack(ErrPhoneIsRegister)
	}

	// 定义用户数据
	userEntity = &models.Users{
		Id:       wuid.GenUid(l.svcCtx.Config.Mysql.DataSource),
		Avatar:   in.Avatar,
		Nickname: in.Nickname,
		Phone:    in.Phone,
		Sex: sql.NullInt64{
			Int64: int64(in.Sex),
			Valid: true,
		},
	}

	if len(in.Password) > 0 {
		genPassword, err := encrypt.GenPasswordHash([]byte(in.Password))
		if err != nil {
			return nil, errors.Wrapf(ErrEncryptionPwd, "encrypt password err %v", err)
		}
		userEntity.Password = sql.NullString{
			String: string(genPassword),
			Valid:  true,
		}
	}

	_, err = l.svcCtx.UsersModel.Insert(l.ctx, userEntity)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "insert user err %v，req %v", err, *userEntity)
	}

	// 生成token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now,
		l.svcCtx.Config.Jwt.AccessExpire, userEntity.Id)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ctxdata get jwt token err %v", err)
	}

	return &user.RegisterResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
