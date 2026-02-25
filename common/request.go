package common

import (
	"errors"
	"strings"
)

type LoginRequest struct {
	Username string `json:"username" `
	Password string `json:"password" `
}
type UpdatePwdRequest struct {
	Username    string `json:"username" binding:"required"`     // 用户名
	OldPassword string `json:"old_password" binding:"required"` // 原密码
	NewPassword string `json:"new_password" binding:"required"` // 新密码
}

// 不返回密码,返回给前端
type UserListResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	UpdatedAt string `json:"updated_at"`
}

// 前端传递请求分页参数
type PageRequest struct {
	Page     int    `form:"page" json:"page" `           // 当前页（默认1）form:"page指定参数的「表单 / URL 查询参数」映射名。
	PageSize int    `form:"page_size" json:"page_size" ` // 每页条数（默认10）
	Keyword  string `form:"keyword" json:"keyword"`      // 搜索关键词（用户名模糊匹配）
}

// 后端返回给前端
type PageResponse struct {
	List       interface{} `json:"list"` // 数据列表
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// 只校验用户输入的格式、非空、基本规则
// 清除空格
func (req *LoginRequest) normalize() {
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

}

func (req *LoginRequest) Validate() error {
	req.normalize()
	if req.Username == "" {
		return errors.New("username is empty")
	}
	if req.Password == "" {
		return errors.New("password is empty")
	}
	return nil
}

// 页码参数校验
func (req *PageRequest) ValidatePageRequest() (string, bool) {
	if req.Page <= 0 {
		return "分页页码必须大于0", false
	}

	if req.PageSize <= 0 {
		return "每页条数必须大于0", false
	}
	if req.PageSize > 100 {
		return "每页条数不能超过100", false
	}

	if len(req.Keyword) > 50 {
		return "搜索关键词长度不能超过50个字符", false
	}

	// 所有校验通过
	return "", true
}
