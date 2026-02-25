package repository

import (
	"Manager/domain/entity"
	"context"
)

// 接口
type UserRepo interface {
	//判断是否用户存在
	IsExist(ctx context.Context, username string) (bool, error)
	//查询数据
	FindById(ctx context.Context, id int) (*entity.Users, error)
	FindByUsername(ctx context.Context, username string) (*entity.Users, error)
	GetUsersList(ctx context.Context, page, pageSize int, keyword string) (list []entity.Users, total int64, err error)
	DeleteUser(ctx context.Context, id int) error

	//操作数据
	UpdatePsd(ctx context.Context, username, newPassword string) error
	// ========== 管理员调用用户模块==========
	AdminUpdateUName(ctx context.Context, username string, newUsername string) error
	AdminUpdateUPsd(ctx context.Context, username string, newPassword string) error
}
