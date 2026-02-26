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

type UserService interface {
	// 登录相关
	UserLogin(ctx context.Context, username, password string) (*entity.Users, error)
	// 密码修改相关
	UpdatePsd(ctx context.Context, username, oldPwd, newPwd string) error
	// 管理员操作相关
	AdminUpdateUName(ctx context.Context, cmd *ChangeUsernameCmd) error
	AdminUpdateUPsd(ctx context.Context, cmd *ChangePasswordCmd) error
	DeleteUser(ctx context.Context, id string) error
	GetUsersList(ctx context.Context, page, pageSize int, keyword string) (resp common.PageResponse, err error)
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
func (u *userService) UpdatePsd(ctx context.Context, username, oldPwd, newPwd string) error {
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
	err = u.userRepo.UpdatePsd(ctx, username, newPwd)
	if err != nil {
		return fmt.Errorf("更新密码失败：%w", err)
	}
	return nil
}

// ------------------------ 管理员操作逻辑 ------------------------
func (u *userService) AdminUpdateUName(ctx context.Context, cmd *ChangeUsernameCmd) error {
	// 调用仓库层更新用户名
	err := u.userRepo.AdminUpdateUName(ctx, cmd.UserID, cmd.NewUsername)
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

func (u *userService) AdminUpdateUPsd(ctx context.Context, cmd *ChangePasswordCmd) error {
	err := u.userRepo.AdminUpdateUPsd(ctx, cmd.UserID, cmd.NewPassword)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("待修改的用户不存在")
		}
		return fmt.Errorf("管理员修改用户密码异常：%w", err)
	}
	return nil
}

func (u *userService) DeleteUser(ctx context.Context, id string) error {
	err := u.userRepo.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("删除的用户不存在")
		}
		return fmt.Errorf("管理员删除用户异常：%w", err)
	}
	return nil
}

// 管理员查询用户列表 接收common.PageRequest中的page，pagesize等参数
func (u *userService) GetUsersList(ctx context.Context, page, pageSize int, keyword string) (resp common.PageResponse, err error) {
	list, total, err := u.userRepo.GetUsersList(ctx, page, pageSize, keyword)
	if err != nil {
		return resp, fmt.Errorf("查询用户列表异常：%w", err)
	}

	var userList []common.UserListResponse //使用不返回密码的
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

	// 计算总页数
	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages += 1
	}

	// 组装分页返回数据
	resp = common.PageResponse{
		List:       userList,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	return resp, nil
}
