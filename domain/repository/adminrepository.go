package repository

import (
	"Manager/domain/entity"
	"context"
)

type AdminRepository interface {
	FindByAdminPassword(ctx context.Context, username, password string) (*entity.Admins, error)
}
