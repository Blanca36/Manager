package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Counter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now().Unix() //会返回Unix 时间戳（纯数字）
		ctx.Next()

		end := time.Now().Unix()
		fmt.Fprintf(ctx.Writer, "接口%s执行时间为：%d秒\n", ctx.Request.URL, end-start)
	}
}

func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		// 执行后续中间件/接口逻辑
		ctx.Next()
	}
}
