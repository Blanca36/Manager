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

// ========================= 用户领域层服务定义=========================

// UserDomain 领域层接口：定义用户核心业务能力
type UserSrv interface {
	UserLogin(ctx context.Context, user *entity.Users) (*entity.Users, error)
	UpdatePsd(ctx context.Context, username, oldPwd, newPwd string) error
	AdminUpdateUName(ctx context.Context, username string, newUsername string) error
	AdminUpdateUPsd(ctx context.Context, username string, newPassword string) error
	DeleteUser(ctx context.Context, id int) error
	GetUsersList(ctx context.Context, page, pageSize int, keyword string) (resp common.PageResponse, err error)
}

// 领域层实现结构体，依赖仓储接口
type userSrv struct {
	userRepo repository.UserRepo
}

// NewUserDomainImpl 领域层构造函数，依赖注入仓储层实例
func NewUserSrv(repo repository.UserRepo) UserSrv {
	return &userSrv{userRepo: repo}
}

// 管理员查询用户列表 接收common.PageRequest中的page，pagesize等参数
func (u *userSrv) GetUsersList(ctx context.Context, page, pageSize int, keyword string) (resp common.PageResponse, err error) {
	list, total, err := u.userRepo.GetUsersList(ctx, page, pageSize, keyword)
	if err != nil {
		return resp, fmt.Errorf("查询用户列表异常：%w", err)
	}
	var userList []common.UserListResponse //使用不返回密码的
	for _, user := range list {
		userList = append(userList, common.UserListResponse{
			ID:        uint(user.ID),
			Username:  user.Username,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"), //提供初始时间参考
			Email:     user.Email,
			Phone:     user.Phone,
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"), // 补充：time.Time 转 string
		})
	}

	// 计算总页数
	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages += 1
	}

	// 4. 组装分页返回数据
	resp = common.PageResponse{
		List:       userList,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	return resp, nil
}

// 用户登录核心逻辑
func (u *userSrv) UserLogin(ctx context.Context, user *entity.Users) (*entity.Users, error) {
	dbUser, err := u.userRepo.FindByUsername(ctx, user.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("登录失败：用户名或密码错误")
		}
		return nil, fmt.Errorf("登录失败：用户信息查询异常，%w", err)
	}

	if dbUser.Password != user.Password {
		return nil, fmt.Errorf("登录失败：用户名或密码错误")
	}
	return dbUser, nil
}

func (u *userSrv) UpdatePsd(ctx context.Context, username, oldPwd, newPwd string) error {
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

// AUpdateUserName 管理员权限：修改普通用户用户名
func (u *userSrv) AdminUpdateUName(ctx context.Context, username string, newUsername string) error {
	// 调用仓库层更新用户名
	err := u.userRepo.AdminUpdateUName(ctx, username, newUsername)
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

// AUpdateUserPsw 管理员权限：重置普通用户密码
func (u *userSrv) AdminUpdateUPsd(ctx context.Context, username string, newPassword string) error {
	err := u.userRepo.AdminUpdateUPsd(ctx, username, newPassword)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("待修改的用户不存在")
		}
		return fmt.Errorf("管理员修改用户密码异常：%w", err)
	}
	return nil
}

func (u *userSrv) DeleteUser(ctx context.Context, id int) error {
	err := u.userRepo.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("删除的用户不存在")
		}
		return fmt.Errorf("管理员删除用户异常：%w", err)
	}
	return nil
}

// ========================= 用户应用层服务=========================

// UserAppService 应用层接口：重命名为专属名称，避免与其他Service冲突，面向前端/外部调用
type UserAppSrv interface {
	UserLogin(ctx context.Context, username, password string) (*entity.Users, error)
	UpdatePsd(ctx context.Context, username, oldPwd, newPwd string) error
}

// UserAppServiceImpl 应用层实现结构体，注入用户领域层服务
type userAppSrv struct {
	userSrv UserSrv
}

// NewUserAppServiceImpl 应用层构造函数，标准化命名，注入领域层依赖
func NewUserAppSvcImpl(srv UserSrv) UserAppSrv {
	return &userAppSrv{userSrv: srv}
}

// LocalLoginUser 应用层登录接口：参数校验 + 调用领域层登录逻辑
func (s *userAppSrv) UserLogin(ctx context.Context, username, password string) (*entity.Users, error) {
	if username == "" || password == "" {
		return nil, fmt.Errorf("用户名和密码不能为空")
	}
	// 组装领域实体，传递给领域层服务
	domainUser := &entity.Users{
		Username: username,
		Password: password,
	}
	userInfo, err := s.userSrv.UserLogin(ctx, domainUser)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

// UpdateUPassword 应用层修改密码接口：参数校验 + 调用领域层修改逻辑
func (s *userAppSrv) UpdatePsd(ctx context.Context, username, oldPwd, newPwd string) error {
	if username == "" || oldPwd == "" || newPwd == "" {
		return fmt.Errorf("用户名、原密码、新密码不能为空")
	}
	err := s.userSrv.UpdatePsd(ctx, username, oldPwd, newPwd)
	if err != nil {
		return err
	}
	return nil
}
