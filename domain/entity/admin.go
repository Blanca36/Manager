package entity

import "time"

type Admins struct {
	ID        int64     `gorm:"primary_key;column:id;table:admins" json:"id"`
	Username  string    `gorm:"type:varchar(50);not null;unique;column:username" json:"username"`
	Password  string    `gorm:"type:varchar(100);not null;column:password" json:"password"`
	Email     string    `gorm:"type:varchar(100);unique;column:email" json:"email"`
	CreatedAt time.Time `gorm:"default:current_timestamp;column:created_at" json:"created_at"` // 创建时间（默认当前时间）
	UpdatedAt time.Time `gorm:"default:current_timestamp;column:updated_at" json:"updated_at"` // 更新时间（默认当前时间）
}

func (a *Admins) TableName() string {
	return "admins"
}
