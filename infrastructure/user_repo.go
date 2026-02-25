package infrastructure

import (
	"Manager/domain/entity"
	"Manager/domain/repository"
	"context"
	"errors"
	"fmt"

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

func (p *PsgUserRepo) FindById(ctx context.Context, id int) (*entity.Users, error) {
	user := &entity.Users{}
	// 注意：GORM的First如果没找到记录，会返回gorm.ErrRecordNotFound错误
	result := p.db.WithContext(ctx).Where("id = ?", id).First(user)
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
func (p *PsgUserRepo) UpdatePsd(ctx context.Context, username, newPassword string) error {
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

// ========== 管理员==========
// 管理员修改用户密码
func (p *PsgUserRepo) AdminUpdateUName(ctx context.Context, username string, newUsername string) error {
	oldExist, err := p.IsExist(ctx, username) //调用isexist
	if err != nil {
		return fmt.Errorf("查询目标用户是否存在异常：%w", err)
	}
	if !oldExist {
		return gorm.ErrRecordNotFound
	}

	// 2. 调用 IsExist 判断新用户名是否已存在
	newExist, err := p.IsExist(ctx, newUsername)
	if err != nil {
		return err // 数据库异常直接返回
	}
	if newExist {
		return gorm.ErrDuplicatedKey // 新用户名重复
	}

	result := p.db.Model(&entity.Users{}).
		Where("username = ?", username).
		Update("username", newUsername)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 管理员修改密码
func (p *PsgUserRepo) AdminUpdateUPsd(ctx context.Context, username string, newPassword string) error {
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

// 获取用户列表
func (p *PsgUserRepo) GetUsersList(ctx context.Context, page, pageSize int, keyword string) (list []entity.Users, total int64, err error) {

	userlist := p.db.WithContext(ctx).Model(&entity.Users{})
	if keyword != "" {
		userlist = userlist.Where("username LIKE ?", "%"+keyword+"%") //模糊查询
	}
	//查询总条数
	err = userlist.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize //标准分页公式
	err = userlist.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (p *PsgUserRepo) DeleteUser(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("用户ID不合法，必须为正整数")
	}
	_, err := p.FindById(ctx, id)
	if err != nil {
		return err
	}
	result := p.db.Where("id = ?", id).Delete(&entity.Users{})
	if result.Error != nil {
		return fmt.Errorf("删除用户失败: %w", result.Error)
	}
	return nil
}
