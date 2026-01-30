package service

import (
	"Manager/domain/entity"
	"Manager/domain/repository"
	"context"
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
)

type UserDomain interface {
	LoginUser(ctx context.Context, user *entity.Users) (*entity.Users, error)
}

type UserDomainImpl struct {
	ur repository.UserRepository //user repo仓库抽象依赖
}

func NewUserDomainImpl(repo repository.UserRepository) UserDomain {
	return &UserDomainImpl{ur: repo}
}

func (u *UserDomainImpl) LoginUser(ctx context.Context, user *entity.Users) (*entity.Users, error) {
	dbUser, err := u.ur.FindByUsernamePassword(ctx, user.Username, user.Password)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("登录失败：用户名或密码错误")
		}
		return nil, fmt.Errorf("登录失败：用户信息查询异常，%w", err)
	}
	if dbUser.Password != user.Password {
		return nil, fmt.Errorf("登录失败：密码错误")
	}
	return dbUser, nil
}
