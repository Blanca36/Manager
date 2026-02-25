package middleware

import (
	"Manager/common"
	"Manager/tool"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			tool.Warn("管理员接口访问-未携带token", zap.String("path", c.FullPath()))
			c.JSON(http.StatusOK, common.UnauthorizedError)
			c.Abort() // 终止请求链，不执行后续接口逻辑
			return
		}

		// 解析token格式（去掉Bearer前缀）
		parts := strings.SplitN(authHeader, " ", 2) //避免 Token 本身包含空格导致分割错误，限制最多分 2 段。
		if len(parts) != 2 || parts[0] != "Bearer" {
			tool.Warn("管理员接口访问-token格式错误", zap.String("auth_header", authHeader))
			c.JSON(http.StatusOK, gin.H{
				"code":    401,
				"message": "token格式错误：请使用Bearer {token}格式",
				"data":    nil,
			})
			c.Abort()
			return
		}
		tokenString := parts[1] // 提取真正的 Token 字符串（去掉 Bearer 前缀）

		// 校验token有效性
		claims, err := tool.ParseAdminToken(tokenString)
		if err != nil {
			tool.Warn("管理员接口访问-token校验失败", zap.Error(err), zap.String("token", tokenString))
			c.JSON(http.StatusOK, gin.H{
				"code":    401,
				"message": "登录态失效：" + err.Error(),
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 将管理员信息存入上下文
		c.Set("admin_id", claims.AdminID)
		c.Set("admin_username", claims.Username)
		c.Next()
	}
}
