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
func (p *PsgUserRepo) AdminUpdateUName(ctx context.Context, id string, newUsername string) error {
	_, err := p.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return fmt.Errorf("查询目标用户是否存在异常：%w", err)
	}
	// 校验新用户名是否已存在
	newExist, err := p.IsExist(ctx, newUsername)
	if err != nil {
		return fmt.Errorf("查询新用户名是否存在异常：%w", err)
	}
	if newExist {
		return gorm.ErrDuplicatedKey // 新用户名重复
	}

	result := p.db.WithContext(ctx).Model(&entity.Users{}).
		Where("id = ?", id).
		Update("username", newUsername)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// 管理员修改密码
func (p *PsgUserRepo) AdminUpdateUPsd(ctx context.Context, id string, newPassword string) error {
	result := p.db.WithContext(ctx).Model(&entity.Users{}).
		Where("id = ?", id).
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

func (p *PsgUserRepo) DeleteUser(ctx context.Context, id string) error {

	if id == "" {
		return errors.New("用户ID不能为空")
	}

	_, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("用户ID格式错误，必须为数字字符串：%w", err)
	}

	_, err = p.FindById(ctx, id)
	if err != nil {
		return fmt.Errorf("待删除的用户不存在：%w", err)
	}

	result := p.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.Users{})
	if result.Error != nil {
		return fmt.Errorf("删除用户失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
