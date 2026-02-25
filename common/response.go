package common

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Success 成功响应
func Success(data interface{}) Response {
	return Response{
		Code:    200,
		Message: "操作成功",
		Data:    data,
	}
}

// ParamError 参数错误响应
func ParamError(msg string) Response {
	return Response{
		Code:    400,
		Message: msg,
		Data:    nil,
	}
}

// ServerError 服务异常响应
func ServerError(msg string) Response {
	return Response{
		Code:    500,
		Message: msg,
		Data:    nil,
	}
}

// UnauthorizedError 鉴权失败响应
func UnauthorizedError(msg string) Response {
	return Response{
		Code:    401,
		Message: msg,
		Data:    nil,
	}
}
