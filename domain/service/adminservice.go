package service

import (
	"Manager/domain/entity"
	"Manager/domain/repository"
	"context"
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
)

type AdminDomain interface {
	LoginAdmin(ctx context.Context, admin *entity.Admins) (*entity.Admins, error)
}
type AdminDomainImpl struct {
	ar repository.AdminRepository //admin repo仓库抽象依赖
}

func NewAdminDomainImpl(repo repository.AdminRepository) AdminDomain {
	return &AdminDomainImpl{ar: repo}
}
func (a AdminDomainImpl) LoginAdmin(ctx context.Context, admin *entity.Admins) (*entity.Admins, error) {
	dbAdmin, err := a.ar.FindByAdminPassword(ctx, admin.Username, admin.Password)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("登录失败：用户名或密码错误")
		}
		return nil, fmt.Errorf("登录失败：用户信息查询异常，%w", err)
	}
	if dbAdmin.Password != admin.Password {
		return nil, fmt.Errorf("登录失败：密码错误")
	}
	return dbAdmin, nil
}
