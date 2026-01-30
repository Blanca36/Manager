package handle

import (
	"Manager/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminService application.Service // 应用层服务接口（
}

func NewAdminHandler(adminService application.Service) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

// LocalLoginUserHandler 登录接口处理函数
func (h *AdminHandler) LocalLoginAdminHandler(c *gin.Context) {
	// 1. 定义请求参数结构体（绑定前端传入的参数）
	type LoginRequest struct {
		Username string `json:"username" binding:"required"` // required：参数校验
		Password string `json:"password" binding:"required"`
	}
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误：" + err.Error(),
			"data":    nil,
		})
		return
	}

	// 调用应用层的
	userInfo, err := h.adminService.LocalLoginUser(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		// 4. 处理业务异常，返回友好响应
		c.JSON(http.StatusOK, gin.H{
			"code":    500,         // 可自定义业务错误码（如401：用户名密码错误）
			"message": err.Error(), // 应用层/领域层返回的错误信息
			"data":    nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "管理员登录成功",
		"data": gin.H{
			"id":       userInfo.ID,
			"username": userInfo.Username,
			// 可扩展：token、昵称、角色等，**务必排除Password字段**
		},
	})
}
