package repository

import (
	"Manager/domain/entity"
	"context"
)

// 接口
type UserRepository interface {
	//判断是否用户存在
	IsExist(ctx context.Context, username string) (bool, error)
	//查询数据
	FindById(ctx context.Context, id int) (*entity.Users, error)
	FindByUsername(ctx context.Context, username string) (*entity.Users, error)
	FindByUsernamePassword(ctx context.Context, username, password string) (*entity.Users, error)

	//操作数据
	GetUser(ctx context.Context, id int64) (*entity.Users, error)
	GetUsers() ([]entity.Users, error)
	UpdateUser(ctx context.Context, user *entity.Users) (*entity.Users, error)
	UpdatePassword(ctx context.Context, id int64, password ...string) (string, error)
}
