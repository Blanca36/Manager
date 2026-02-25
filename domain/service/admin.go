package service

import (
	"Manager/common"
	"Manager/domain/entity"
	"Manager/domain/repository"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// ========================= 应用层服务=========================
type AdminAppSrv interface {
	AdminLogin(ctx context.Context, username, password string) (*entity.Admins, error)
	AdminUpdateUName(ctx context.Context, username string, newUsername string) error
	AdminUpdateUPsd(ctx context.Context, username string, newPassword string) error
	GetUsersList(ctx context.Context, req common.PageRequest) (resp common.PageResponse, err error)
	DeleteUser(ctx context.Context, id int) error
}

// AdminServiceImpl 应用层实现结构体，注入领域层服务依赖
type AdminSrvImpl struct {
	adminRepo  repository.AdminRepo
	userDomain UserSrv // 用户领域层服务（保留原有依赖）
}

// NewAdminServiceImpl 应用层构造函数，修复原命名笔误 Ne -> New
func NewAdminAppSrvImpl(repo repository.AdminRepo, userDomain UserSrv) AdminAppSrv {
	return &AdminSrvImpl{
		adminRepo:  repo,
		userDomain: userDomain,
	}
}

// 前端通过Handle请求传递common.PageRequest进来，返回 common.PageResponse
func (s *AdminSrvImpl) GetUsersList(ctx context.Context, req common.PageRequest) (resp common.PageResponse, err error) {
	resp, err = s.userDomain.GetUsersList(ctx, req.Page, req.PageSize, req.Keyword)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// LocalLoginAdmin 应用层登录接口：参数校验 + 调用领域层登录逻辑
func (s *AdminSrvImpl) AdminLogin(ctx context.Context, username, password string) (*entity.Admins, error) {
	if username == "" || password == "" {
		return nil, fmt.Errorf("用户名和密码不能为空")
	}
	dbAdmin, err := s.adminRepo.FindByAdminName(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("登录失败：用户名或密码错误")
		}
		return nil, fmt.Errorf("登录失败：用户信息查询异常，%w", err)
	}
	if dbAdmin.Password != dbAdmin.Password {
		return nil, fmt.Errorf("登录失败：密码错误")
	}
	return dbAdmin, nil
}

// AUpdateUserName 应用层接口：管理员修改普通用户名
func (s *AdminSrvImpl) AdminUpdateUName(ctx context.Context, username string, newUsername string) error {
	if username == "" || newUsername == "" {
		return fmt.Errorf("原用户名、新用户名不能为空")
	}
	// 调用用户领域层方法
	err := s.userDomain.AdminUpdateUName(ctx, username, newUsername)
	if err != nil {
		return err
	}
	return nil
}

// AUpdateUserPsw 应用层接口：管理员修改普通用户密码
func (s *AdminSrvImpl) AdminUpdateUPsd(ctx context.Context, username string, newPassword string) error {
	if username == "" || newPassword == "" {
		return fmt.Errorf("待修改的用户名、新密码不能为空")
	}
	// 调用用户领域层方法
	err := s.userDomain.AdminUpdateUPsd(ctx, username, newPassword)
	if err != nil {
		return err
	}
	return nil
}

func (s *AdminSrvImpl) DeleteUser(ctx context.Context, id int) error {
	err := s.userDomain.DeleteUser(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
