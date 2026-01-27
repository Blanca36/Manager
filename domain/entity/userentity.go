package entity

// 面向前端
type User struct {
	Id         int64  `gorm:"column:id;primary_key;not null;comment:唯一id"`
	Username   string `gorm:"column:username;type:varchar(64);index:idx_tbl_users_username;not null;comment:用户名'"`
	Password   string `gorm:"column:password;type:varchar(64);not null;comment:用户名'" `
	CreateTime int64  `gorm:"column:create_time;type:bigint;index:idx_tbl_users_create_time;not null "`
	UpdateTime int64  `gorm:"column:update_time;type:bigint;index:idx_tbl_users_update_time;not null "`
}

func (User) TableName() string {
	return "tbl_users"
}
