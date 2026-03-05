package service

import (
	"Manager/domain/entity"
	"Manager/domain/repository"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type UserService interface {
	// 登录相关
	UserLogin(ctx context.Context, username, password string) (*entity.Users, error)
	// 密码修改相关
	UpdatePassword(ctx context.Context, username, oldPwd, newPwd string) error
}

// 实现结构体，依赖仓储接口
type userService struct {
	userRepo repository.UserRepo
}

// 构造函数，依赖注入仓储层实例
func NewUserService(repo repository.UserRepo) UserService {
	return &userService{userRepo: repo}
}

// ------------------------ 登录 -----------------------
func (u *userService) UserLogin(ctx context.Context, username, password string) (*entity.Users, error) {
	// 1. 原应用层的参数校验逻辑
	if username == "" || password == "" {
		return nil, fmt.Errorf("用户名和密码不能为空")
	}
	dbUser, err := u.userRepo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("登录失败：用户名或密码错误")
		}
		return nil, fmt.Errorf("登录失败：用户信息查询异常，%w", err)
	}

	if dbUser.Password != password {
		return nil, fmt.Errorf("登录失败：用户名或密码错误")
	}
	return dbUser, nil
}

// ------------------------ 修改密码 ------------------------
func (u *userService) UpdatePassword(ctx context.Context, username, oldPwd, newPwd string) error {
	if username == "" || oldPwd == "" || newPwd == "" {
		return fmt.Errorf("用户名、原密码、新密码不能为空")
	}

	result, err := u.userRepo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("用户名不存在，修改失败")
		}
		return fmt.Errorf("查询用户信息异常：%w", err)
	}
	// 对比旧密码
	if result.Password != oldPwd {
		return fmt.Errorf("原密码错误，修改失败")
	}
	err = u.userRepo.UpdatePassword(ctx, username, newPwd)
	if err != nil {
		return fmt.Errorf("更新密码失败：%w", err)
	}
	return nil
}
