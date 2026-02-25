package handler

import (
	"Manager/common"
	"Manager/domain/service"
	"Manager/tool"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService service.UserAppSrv // 应用层服务接口（
}

func NewUserHandler(userService service.UserAppSrv) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// LocalLoginUserHandler 登录接口处理函数
func (h *UserHandler) UserLogin(c *gin.Context) {
	// 1. 定义请求参数结构体（绑定前端传入的参数）

	req := &common.LoginRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		// 记录参数错误日志
		tool.Error("用户登录-参数绑定失败",
			zap.String("username", req.Username),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, common.ParamError("请求参数错误："+err.Error()))
		return
	}
	tool.Info("用户登录-接收请求",
		zap.String("username", req.Username),
	)

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, common.ParamError("请求参数错误："+err.Error()))
		return
	}
	// 调用应用层
	userInfo, err := h.userService.UserLogin(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		tool.Error("用户登录-业务处理失败",
			zap.String("username", req.Username),
			zap.String("password", req.Password),
			zap.Error(err),
		)
		c.JSON(http.StatusUnauthorized, common.UnauthorizedError(err.Error()))
		return
	}
	reqID := c.GetString("X-Request-ID")
	tool.Info("用户登录-成功",
		zap.String("req_id", reqID),
		zap.String("username", req.Username),
		zap.Uint("user_id", uint(userInfo.ID)),
	)

	data := gin.H{
		"id":       userInfo.ID,
		"username": userInfo.Username,
	}
	c.JSON(http.StatusOK, common.Success(data))
}

// =====================用户修改密码=====================
func (h *UserHandler) UpdatePsd(c *gin.Context) {

	var req common.UpdatePwdRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		tool.Error("用户密码参数绑定失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ParamError("请求参数错误："+err.Error()))
		return
	}
	//改密记录
	reqID := c.GetString("X-Request-ID")
	tool.Info("用户改密-接收请求",
		zap.String("req_id", reqID),
		zap.String("username", req.Username),
	)

	err := h.userService.UpdatePsd(
		c.Request.Context(),
		req.Username,
		req.OldPassword,
		req.NewPassword,
	)
	if err != nil {
		tool.Error("用户改密-业务处理失败",
			zap.String("req_id", reqID),
			zap.String("username", req.Username),
			zap.Error(err),
		)
		c.JSON(http.StatusOK, common.ServerError(err.Error()))
		return
	}

	tool.Info("用户改密-成功",
		zap.String("req_id", reqID),
		zap.String("username", req.Username),
	)
	// 无数据返回，传入nil即可
	c.JSON(http.StatusOK, common.Success(nil))
}

func ErrHandler(c *gin.Context, err error) {
	var myErr MyError
	if errors.As(err, &myErr) {
		slog.Info("err", slog.Any("err", myErr.Case()), slog.String("requestId", c.GetString("X-Request-ID")))
		c.JSON(myErr.Code, myErr.Message)
		return
	} else {
		c.JSON(http.StatusInternalServerError, "ccdcd")
		return
	}
}

func MyErrorWith(code int, msg string) *MyError {
	return &MyError{
		Code:    code,
		Message: msg,
	}
}

type MyError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	err     error  `json:"-"`
}

func (e *MyError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

func (e *MyError) WithCase(err error) *MyError {
	var ne *MyError
	safe.Copy(&e, e)
	ne.err = err
	return ne
}

func (e *MyError) Case() error {
	return e.err
}
