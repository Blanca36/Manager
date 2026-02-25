package handler

import (
	"Manager/common"
	"Manager/domain/service"
	"Manager/tool"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AdminHandler struct {
	adminSrv service.AdminAppSrv // 应用层服务接口（
}

func NewAdminHandler(adminSrv service.AdminAppSrv) *AdminHandler {
	return &AdminHandler{
		adminSrv: adminSrv,
	}
}

// ========================管理员登录====================================
// LocalLoginUserHandler 登录接口处理函数
func (h *AdminHandler) AdminLogin(c *gin.Context) {
	// 1. 定义请求参数结构体（绑定前端传入的参数）

	var req common.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		// 记录参数错误日志
		tool.Error("管理员登录-参数绑定失败",
			zap.String("username", req.Username),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, common.ParamError("请求参数错误："+err.Error()))
		return
	}
	if err := (&req).Validate(); err != nil {
		c.JSON(http.StatusBadRequest, common.ParamError("请求参数错误："+err.Error()))
		return
	}

	// 调用应用层的
	adminInfo, err := h.adminSrv.AdminLogin(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		// 4. 处理业务异常，返回友好响应
		c.JSON(http.StatusUnauthorized, common.UnauthorizedError(err.Error()))
		return

	}
	token, err := tool.GenerateAdminToken(uint(adminInfo.ID), adminInfo.Username)
	if err != nil {
		tool.Error("管理员登录-生成token失败", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "登录失败：生成token异常",
			"data":    nil,
		})
		return
	}

	// 4. 返回成功结果（treqID := c.GetString("X-Request-ID")oken需前端存储，后续请求携带）
	reqID := c.GetString("X-Request-ID")
	tool.Info("管理员登录成功",
		zap.String("req_id", reqID),
		zap.String("username", req.Username),
		zap.Uint("admin_id", uint(adminInfo.ID)),
		zap.String("username", adminInfo.Username))

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data": gin.H{
			"token":  token,     // 返回token
			"expire": "2小时",     // 提示过期时间
			"admin":  adminInfo, // 返回管理员基本信息（非敏感）
		},
	})
}

// ========================获取用户列表====================================
func (h *AdminHandler) GetUsersList(c *gin.Context) {
	var req common.PageRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		tool.Error("管理员查询用户列表-参数绑定失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ParamError("请求参数错误："+err.Error()))
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if msg, ok := (&req).ValidatePageRequest(); !ok {
		c.JSON(http.StatusBadRequest, common.ParamError(msg))
		return
	}

	reqID := c.GetString("X-Request-ID")
	tool.Info("管理员查询用户列表-接收请求",
		zap.String("req_id", reqID),
		zap.Int("page", req.Page),
		zap.Int("page_size", req.PageSize),
		zap.String("keyword", req.Keyword),
	)

	resp, err := h.adminSrv.GetUsersList(c.Request.Context(), req)
	if err != nil {
		tool.Error("管理员查询用户列表-业务处理失败",
			zap.String("req_id", reqID),
			zap.Error(err),
		)
		c.JSON(http.StatusOK, common.ServerError(err.Error()))
		return
	}

	// 5. 响应成功
	tool.Info("管理员查询用户列表-成功",
		zap.String("req_id", reqID),
		zap.Int64("total", resp.Total),
		zap.Int("total_pages", resp.TotalPages),
	)
	c.JSON(http.StatusOK, common.Success(resp.List))
}

// ========================管理员修改用户名====================================
func (h *AdminHandler) AdminUpdateUName(c *gin.Context) {

	type UpdateUsernameRequest struct {
		OldUsername string `json:"old_username" binding:"required"` // 原用户名
		NewUsername string `json:"new_username" binding:"required"` // 新用户名
	}

	var req UpdateUsernameRequest
	// 2. 绑定参数并校验
	if err := c.ShouldBindJSON(&req); err != nil {
		tool.Error("管理员修改用户用户名-参数绑定失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ParamError("请求参数错误："+err.Error()))
		return
	}

	// 3. 记录日志
	reqID := c.GetString("X-Request-ID")
	tool.Info("管理员修改用户用户名-接收请求",
		zap.String("req_id", reqID),
		zap.String("old_username", req.OldUsername),
		zap.String("new_username", req.NewUsername),
	)

	// 4. 调用应用层方法
	err := h.adminSrv.AdminUpdateUName(
		c.Request.Context(),
		req.OldUsername,
		req.NewUsername,
	)
	if err != nil {
		tool.Error("管理员修改用户用户名-业务处理失败",
			zap.String("req_id", reqID),
			zap.String("old_username", req.OldUsername),
			zap.Error(err),
		)
		c.JSON(http.StatusOK, common.ServerError(err.Error()))
		return
	}

	tool.Info("管理员修改用户用户名-成功",
		zap.String("req_id", reqID),
		zap.String("old_username", req.OldUsername),
		zap.String("new_username", req.NewUsername),
	)
	c.JSON(http.StatusOK, common.Success(nil))
}

// ========================管理员修改用户密码====================================
func (h *AdminHandler) AdminUpdateUPsd(c *gin.Context) {
	//id := c.Param("id")

	var req common.UpdatePwdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tool.Error("管理员修改用户密码-参数绑定失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ParamError("请求参数错误："+err.Error()))
		return
	}
	reqID := c.GetString("X-Request-ID")
	tool.Info("管理员修改用户密码-接收请求",
		zap.String("req_id", reqID),
		zap.String("target_username", req.Username),
		zap.String("new_password", req.NewPassword),
	)
	err := h.adminSrv.AdminUpdateUPsd(
		c.Request.Context(),
		req.Username,
		req.NewPassword,
	)
	if err != nil {
		tool.Error("管理员修改用户密码-业务处理失败",
			zap.String("req_id", reqID),
			zap.String("target_username", req.Username),
			zap.Error(err),
		)
		c.JSON(http.StatusOK, common.ServerError(err.Error()))
		return

	}

	// 5. 响应成功
	tool.Info("管理员修改用户密码-成功",
		zap.String("req_id", reqID),
		zap.String("target_username", req.Username),
	)
	c.JSON(http.StatusOK, common.Success(nil))
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	type DeleteUserRequest struct {
		UserID uint `json:"user_id"` //笔误用了小写userID，导致id绑定失败
	}
	var req DeleteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tool.Error("管理员删除用户-参数绑定失败",
			zap.Error(err),
			zap.String("req_id", c.GetString("X-Request-ID")))
		c.JSON(http.StatusBadRequest, common.ParamError("请求参数错误："+err.Error()))
		return // 关键：添加return，阻止后续代码执行,缺少了return导致用id=0调用后续代码
	}
	//// 2. 记录日志（操作人+要删除的用户ID）
	//adminID := c.GetUint("admin_id")               // 从上下文获取登录的管理员ID
	//adminUsername := c.GetString("admin_username") // 从上下文获取管理员用户名
	//tool.Info("管理员发起删除用户请求",
	//	zap.Uint("operator_admin_id", adminID),
	//	zap.String("operator_username", adminUsername),
	//	zap.Uint("delete_user_id", req.UserID),
	//)

	if req.UserID <= 0 {
		tool.Error("管理员删除用户-用户ID不合法",
			zap.Uint("delete_user_id", req.UserID),
			zap.String("req_id", c.GetString("X-Request-ID")))
		c.JSON(http.StatusBadRequest, common.ParamError("用户ID必须为正整数"))
		return
	}

	err := h.adminSrv.DeleteUser(c.Request.Context(), int(req.UserID))
	if err != nil {
		tool.Error("管理员删除失败",
			zap.Uint("delete_user_id", req.UserID),
			zap.Error(err),
			zap.String("req_id", c.GetString("X-Request-ID")))
		c.JSON(http.StatusOK, common.ServerError(err.Error()))
		return
	}
	tool.Info("管理员删除用户成功",
		zap.Uint("delete_user_id", req.UserID),
		zap.String("req_id", c.GetString("X-Request-ID")),
	)
	c.JSON(http.StatusOK, common.Success(nil))
}
