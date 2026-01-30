package infrastructure

import (
	"Manager/domain/entity"
	"Manager/domain/repository"
	"context"
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
)

// 数据持久化（Domain 层 repository 的具体实现，数据库 CRUD 操作）
type PsgUserRepository struct {
	db *gorm.DB
}

func NewPsgUserRepository(db *gorm.DB) repository.UserRepository {
	return &PsgUserRepository{db: db}
}

func (p PsgUserRepository) wrapErr(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("用户不存在")
	} else {
		return errors.New("内部错误")
	}
}
func (p PsgUserRepository) IsExist(ctx context.Context, username string) (bool, error) {
	userPO := &entity.Users{}
	result := p.db.Where("username = ?", username).First(userPO)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) { //errors.Is专门用于判断一个错误是否为指定的目标错误
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}

func (p PsgUserRepository) FindById(ctx context.Context, id int) (*entity.Users, error) {
	user := &entity.Users{}
	result := p.db.Where("id = ?", id).First(user)
	if result.Error != nil {
		return nil, fmt.Errorf("查询id失败: %w", id, result.Error)
	}
	return user, nil
}

func (p PsgUserRepository) FindByUsername(ctx context.Context, username string) (*entity.Users, error) {
	user := &entity.Users{}
	result := p.db.Where("username = ?", username).First(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在[username:%s]", username)
		}
		return nil, fmt.Errorf("查询用户失败: %w", username, result.Error)
	}
	return user, nil
}

func (p PsgUserRepository) FindByUsernamePassword(ctx context.Context, username, password string) (*entity.Users, error) {
	user := &entity.Users{}
	result := p.db.Where("username = ? AND password = ?", username, password).First(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (p PsgUserRepository) UpdateUser(ctx context.Context, user *entity.Users) (*entity.Users, error) {

	panic("implement me")
}

func (p PsgUserRepository) GetUser(ctx context.Context, id int64) (*entity.Users, error) {

	panic("implement me")
}

// 获取所有用户列表
func (p PsgUserRepository) GetUsers() ([]entity.Users, error) {
	panic("implement me")
}

func (p PsgUserRepository) UpdatePassword(ctx context.Context, id int64, password ...string) (string, error) {

	panic("implement me")
}

func (p PsgUserRepository) GetByName(ctx context.Context, username string) (*entity.Users, error) {

	panic("implement me")

}
