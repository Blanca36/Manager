package tool

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 全局logger实例
var Logger *zap.Logger

// InitLogger 初始化日志配置
func InitLogger() {
	// 开发环境配置（控制台输出，格式友好）
	config := zap.NewDevelopmentConfig()
	// 调整日志级别（Debug/Info/Warn/Error，级别越低输出越多）
	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	// 禁用调用栈（开发环境可开启，生产环境关闭）
	config.Development = true

	// 初始化logger
	var err error
	Logger, err = config.Build()
	if err != nil {
		panic("日志初始化失败：" + err.Error())
	}

	defer Logger.Sync()
}

func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

func Warn(s string, fields ...zap.Field) {
	Logger.Warn(s, fields...)
}
