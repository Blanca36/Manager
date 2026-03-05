package repository

import (
	"Manager/domain/entity"
	"context"
)

type AdminRepo interface {
	//GetUsers() ([]entity.Users, error) //获取user列表
	FindByAdminName(ctx context.Context, username string) (*entity.Admins, error)

	//操作用户
	GetUsersList(ctx context.Context, page, pageSize int, keyword string) (list []entity.Users, total int64, err error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUsername(ctx context.Context, id string, newUsername string) error
	UpdateUserPassword(ctx context.Context, id string, newPassword string) error
}

//type Manager interface {
//	Transaction(ctx context.Context, fn func(ctx context.Context, tx Manager) error) error  数据库事务
//	GetUserRepo() UserRepository
//  GetOrderRepo() OderRepository
//}
//// 原子性
//type UserRepository interface {    UserRepository 的设计目标是只负责用户数据的 CRUD 操作（数据访问层）
//	Create(ctx context.Context, user *entity.User) error
//	// GetByName 查询用户
//
//}

//事务往往需要跨多个仓储操作，若有其他业务如OderRepository，则可以单独写
//type OderRepository interface {
//	Create(ctx context.Context, user *entity.User) error
//}  相当于如果有多个Repository就将所有Repository放进Manger中，需要的时候直接调用Manger
