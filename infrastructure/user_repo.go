package infrastructure

import (
	"Manager/domain/entity"
	"Manager/domain/repository"
	"context"
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

// 数据持久化（Domain 层 repository 的具体实现，数据库 CRUD 操作）
type PsgUserRepo struct {
	db *gorm.DB
}

func NewPsgUserRepo(db *gorm.DB) repository.UserRepo {
	return &PsgUserRepo{db: db}
}

func (p *PsgUserRepo) wrapErr(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("用户不存在")
	} else if errors.Is(err, gorm.ErrDuplicatedKey) {
		return errors.New("user is Existed")
	} else {
		return errors.New("内部错误")
	}
}
func (p *PsgUserRepo) IsExist(ctx context.Context, username string) (bool, error) {
	if username == "" {
		return false, fmt.Errorf("用户名不能为空")
	}
	//无需查询整条用户记录，只需统计数量）
	var count int64
	result := p.db.WithContext(ctx).
		Model(&entity.Users{}).
		Where("username = ?", username).
		Count(&count)
	if result.Error != nil {

		return false, fmt.Errorf("查询用户名是否存在异常：%w", result.Error)
	}
	// 数量>0 说明存在，否则不存在
	return count > 0, nil
}

func (p *PsgUserRepo) FindById(ctx context.Context, id string) (*entity.Users, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("用户ID格式错误：%w", err)
	}
	user := &entity.Users{}
	// 注意：GORM的First如果没找到记录，会返回gorm.ErrRecordNotFound错误
	result := p.db.WithContext(ctx).Where("id = ?", idInt).First(user)
	if result.Error != nil {
		// 修正：%w 只绑定错误，id 用 %d 展示（便于日志排查）
		return nil, fmt.Errorf("查询id为%d的用户失败: %w", id, result.Error)
	}
	return user, nil
}

func (p *PsgUserRepo) FindByUsername(ctx context.Context, username string) (*entity.Users, error) {
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

// 只需要返回修改错误error
func (p *PsgUserRepo) UpdatePassword(ctx context.Context, username, newPassword string) error {
	result := p.db.WithContext(ctx).Model(&entity.Users{}).
		Where("username = ?", username).
		Update("password", newPassword)
	if result.Error != nil {
		return result.Error
	}
	//检查受影响的行数，如果行数为 0,则返回为错误
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (p *PsgUserRepo) GetByName(ctx context.Context, username string) (*entity.Users, error) {

	panic("implement me")

}
