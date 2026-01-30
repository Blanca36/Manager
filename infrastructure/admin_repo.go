package infrastructure

import (
	"Manager/domain/entity"
	"Manager/domain/repository"
	"context"

	"github.com/jinzhu/gorm"
)

type PsgAdminRepository struct {
	db *gorm.DB
}

func NewPsgAdminRepository(db *gorm.DB) repository.AdminRepository {
	return &PsgAdminRepository{db: db}
}

func (p PsgAdminRepository) FindByAdminPassword(ctx context.Context, username, password string) (*entity.Admins, error) {
	admin := &entity.Admins{}
	result := p.db.Where("username = ? AND password = ?", username, password).First(admin)
	if result.Error != nil {
		return nil, result.Error
	}
	return admin, nil
}
