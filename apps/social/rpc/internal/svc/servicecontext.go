package svc

import (
	"easy-chat/apps/social/rpc/internal/config"
	"easy-chat/apps/social/socialmodels"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	socialmodels.FriendRequestsModel
	socialmodels.FriendsModel
	socialmodels.GroupMembersModel
	socialmodels.GroupRequestsModel
	socialmodels.GroupsModel
}

func NewServiceContext(c config.Config) *ServiceContext {

	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config: c,

		FriendsModel:        socialmodels.NewFriendsModel(sqlConn, c.Cache),
		FriendRequestsModel: socialmodels.NewFriendRequestsModel(sqlConn, c.Cache),
		GroupsModel:         socialmodels.NewGroupsModel(sqlConn, c.Cache),
		GroupRequestsModel:  socialmodels.NewGroupRequestsModel(sqlConn, c.Cache),
		GroupMembersModel:   socialmodels.NewGroupMembersModel(sqlConn, c.Cache),
	}
}
