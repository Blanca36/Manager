package infrastructure

import (
	"Manager/domain/entity"
	"Manager/domain/repository"
	"context"
	"errors"
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
	p.db.Where("username = ?", username).First(userPO)
	return userPO != nil, nil
}

func (p PsgUserRepository) FindById(ctx context.Context, id int) (*entity.Users, error) {
	panic("implement me")
}

func (p PsgUserRepository) FindByName(ctx context.Context, name string) (*entity.Users, error) {
	panic("implement me")
}

func (p PsgUserRepository) FindByUsername(ctx context.Context, username string) (*entity.Users, error) {
	panic("implement me")
}

func (p PsgUserRepository) FindByUsernamePassword(ctx context.Context, username, password string) (*entity.Users, error) {
	panic("implement me")
}

func (p PsgUserRepository) SaveUser(ctx context.Context, user *entity.Users) (*entity.Users, error) {

	panic("implement me")
}

func (p PsgUserRepository) DeleteUser(ctx context.Context, user *entity.Users) (*entity.Users, error) {
	panic("implement me")
}

func (p PsgUserRepository) UpdateUser(ctx context.Context, user *entity.Users) (*entity.Users, error) {

	panic("implement me")
}

func (p PsgUserRepository) GetUser(ctx context.Context, id int64) (*entity.Users, error) {

	panic("implement me")
}

func (p PsgUserRepository) UpdatePassword(ctx context.Context, id int64, password ...string) (string, error) {

	panic("implement me")
}

func (p PsgUserRepository) GetByName(ctx context.Context, username string) (*entity.Users, error) {

	panic("implement me")

}
