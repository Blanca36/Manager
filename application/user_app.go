package application

import (
	"Manager/domain/entity"
	"Manager/domain/service"
	"context"
	"fmt"
)

// 面向前端业务能力接口
type Service interface {
	LocalLoginUser(ctx context.Context, username, password string) (*entity.Users, error)
}

// 领域层的用户领域服务接口
type ServiceImpl struct {
	ud service.UserDomain //userdomain
}

func NewServiceImpl(srv service.UserDomain) Service {
	return &ServiceImpl{ud: srv}
}
func (s ServiceImpl) LocalLoginUser(ctx context.Context, username, password string) (*entity.Users, error) {
	if username == "" || password == "" {
		return nil, fmt.Errorf("用户名和密码不能为空")
	}
	domainUser := &entity.Users{
		Username: username,
		Password: password,
	}
	userInfo, err := s.ud.LoginUser(ctx, domainUser)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}
