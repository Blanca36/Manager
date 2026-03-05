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

type PsgAdminRepo struct {
	db       *gorm.DB
	userRepo repository.UserRepo // 组合基础 Repo，复用查询逻辑
}

func NewPsgAdminRepo(db *gorm.DB, userRepo repository.UserRepo) repository.AdminRepo {
	return &PsgAdminRepo{
		db:       db,
		userRepo: userRepo,
	}
}

func (p *PsgAdminRepo) FindByAdminName(ctx context.Context, username string) (*entity.Admins, error) {
	admin := &entity.Admins{}
	result := p.db.WithContext(ctx).Where("username = ? ", username).First(admin)
	err := result.Error
	if err != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("管理员不存在[username:%s]", username)
		}
		return nil, fmt.Errorf("查询管理员失败: %w", username, result.Error)
	}
	return admin, nil
}

func (p *PsgAdminRepo) GetUsersList(ctx context.Context, page, pageSize int, keyword string) (list []entity.Users, total int64, err error) {
	userlist := p.db.WithContext(ctx).Model(&entity.Users{})
	if keyword != "" {
		userlist = userlist.Where("username like ?", "%"+keyword+"%") //模糊查询
	}
	//查询总数
	err = userlist.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err = userlist.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (p *PsgAdminRepo) DeleteUser(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("用户ID不能为空")
	}
	_, err := strconv.Atoi(id) //Atoi将字符串类型转换为整数类型
	if err != nil {
		return fmt.Errorf("用户ID格式错误，必须为数字字符串：%w", err)
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

func (p *PsgAdminRepo) UpdateUsername(ctx context.Context, id string, newUsername string) error {
	_, err := p.userRepo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		//标准库 errors.Is 函数判断捕获的错误是否为 GORM 框架定义的「记录未找到」错误（gorm.ErrRecordNotFound），.
		//若是则返回该明确的错误类型，便于上层代码精准处理数据库查询无匹配记录的场景。
		return fmt.Errorf("查询目标用户是否存在异常：%w", err)
	}
	// 校验新用户名是否已存在
	newExist, err := p.userRepo.IsExist(ctx, newUsername)
	if err != nil {
		return fmt.Errorf("查询新用户名是否存在异常：%w", err)
	}
	if newExist {
		return gorm.ErrDuplicatedKey // 新用户名重复，数据库唯一键 / 主键冲突错误
	}
	result := p.db.WithContext(ctx).Model(&entity.Users{}).Where("id = ?", id).Update("username", newUsername)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (p *PsgAdminRepo) UpdateUserPassword(ctx context.Context, id string, newPassword string) error {
	result := p.db.WithContext(ctx).Model(&entity.Users{}).Where("id = ?", id).Update("password", newPassword)
	if result.Error != nil {
		return result.Error
	}
	//检查受影响的行数，如果行数为 0,则返回为错误
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

//func (p *PsgAdminRepo) GetUsers() ([]entity.Users, error) {
//
//	panic("implement me")
//}
