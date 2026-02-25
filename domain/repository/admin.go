package repository

import (
	"Manager/domain/entity"
	"context"
)

type AdminRepo interface {
	GetUsers() ([]entity.Users, error) //获取user列表
	FindByAdminName(ctx context.Context, username string) (*entity.Admins, error)
}
