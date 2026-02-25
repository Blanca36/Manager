package infrastructure

import (
	"Manager/domain/entity"
	"Manager/domain/repository"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type PsgAdminRepo struct {
	db *gorm.DB
}

func NewPsgAdminRepo(db *gorm.DB) repository.AdminRepo {
	return &PsgAdminRepo{db: db}
}

func (p *PsgAdminRepo) FindByAdminName(ctx context.Context, username string) (*entity.Admins, error) {
	admin := &entity.Admins{}
	result := p.db.Where("username = ? ", username).First(admin)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("管理员不存在[username:%s]", username)
		}
		return nil, fmt.Errorf("查询管理员失败: %w", username, result.Error)
	}
	return admin, nil
}

func (p *PsgAdminRepo) GetUsers() ([]entity.Users, error) {

	panic("implement me")
}
