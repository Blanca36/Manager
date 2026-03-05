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

// 管理员修改用户密码的命令结构体
type ChangePasswordCmd struct {
	UserID      string `json:"user_id"`                                      // 目标用户ID
	NewPassword string `json:"new_password" binding:"required,min=6,max=20"` // 新密码
}

// ChangeUsernameCmd 管理员修改用户用户名的命令结构体
type ChangeUsernameCmd struct {
	UserID      string `json:"user_id"`
	NewUsername string `json:"new_username" `
}

type AdminService interface {
	AdminLogin(ctx context.Context, username, password string) (*entity.Admins, error)
	UpdateUsername(ctx context.Context, cmd *ChangeUsernameCmd) error
	UpdateUserPassword(ctx context.Context, cmd *ChangePasswordCmd) error
	GetUsersList(ctx context.Context, req common.PageRequest) (resp common.PageResponse, err error)
	DeleteUser(ctx context.Context, id string) error
}

// AdminServiceImpl 应用层实现结构体，注入领域层服务依赖
type adminSerivce struct {
	adminRepo repository.AdminRepo
	//userService UserService // 用户领域层服务（保留原有依赖）
}

// NewAdminServiceImpl 应用层构造函数
func NewAdminServiceImpl(repo repository.AdminRepo) AdminService {
	return &adminSerivce{
		adminRepo: repo,
		//userService: userService,
	}
}

// 前端通过Handle请求传递common.PageRequest进来，返回 common.PageResponse
func (s *adminSerivce) GetUsersList(ctx context.Context, req common.PageRequest) (resp common.PageResponse, err error) {
	list, total, err := s.adminRepo.GetUsersList(ctx, req.Page, req.PageSize, req.Keyword)
	if err != nil {
		return resp, fmt.Errorf("查询用户列表异常：%w", err)
	}

	var userList []common.UserListResponse
	for _, user := range list {
		userList = append(userList, common.UserListResponse{
			ID:        uint(user.ID),
			Username:  user.Username,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			Email:     user.Email,
			Phone:     user.Phone,
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize != 0 {
		totalPages += 1
	}

	resp = common.PageResponse{
		List:       userList,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}

	return resp, nil
}

// LocalLoginAdmin 应用层登录接口：参数校验 + 调用领域层登录逻辑
func (s *adminSerivce) AdminLogin(ctx context.Context, username, password string) (*entity.Admins, error) {
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
	if dbAdmin.Password != password {
		return nil, fmt.Errorf("登录失败：密码错误")
	}
	return dbAdmin, nil
}

// AUpdateUserName 应用层接口：管理员修改普通用户名
func (s *adminSerivce) UpdateUsername(ctx context.Context, cmd *ChangeUsernameCmd) error {
	if cmd.UserID == "" || cmd.NewUsername == "" {
		return fmt.Errorf("用户ID、新用户名不能为空") // 校验ID而非原用户名
	}
	// 调用用户领域层方法
	err := s.adminRepo.UpdateUsername(ctx, cmd.UserID, cmd.NewUsername)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("待修改的用户不存在")
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fmt.Errorf("新用户名已存在，修改失败")
		}
		return fmt.Errorf("管理员修改用户用户名异常：%w", err)

	}
	return nil
}

// AUpdateUserPsw 应用层接口：管理员修改普通用户密码 UpdateUsername
func (s *adminSerivce) UpdateUserPassword(ctx context.Context, cmd *ChangePasswordCmd) error {
	if cmd.UserID == "" || cmd.NewPassword == "" {
		return fmt.Errorf("待修改的用户ID、新密码不能为空")
	}
	// 调用用户领域层方法
	err := s.adminRepo.UpdateUserPassword(ctx, cmd.UserID, cmd.NewPassword) //调用cmd changepossword
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("待修改的用户不存在")
		}
		return fmt.Errorf("管理员修改用户密码异常：%w", err)
	}
	return nil
}

func (s *adminSerivce) DeleteUser(ctx context.Context, id string) error {
	err := s.adminRepo.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("删除的用户不存在")
		}
		return fmt.Errorf("管理员删除用户异常：%w", err)
	}
	return nil
}
