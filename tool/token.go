package tool

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var adminJwtSecret = []byte("de49a1dc68cbb671e13a8df2deab5f7b") //在控制台使用openssl rand -hex 16生成

type AdminClaims struct {
	AdminID              uint   `json:"admin_id"`
	Username             string `json:"username"`
	jwt.RegisteredClaims        // 内置标准声明（过期时间、签发时间等）
}

func GenerateAdminToken(adminID uint, username string) (string, error) {
	expireTime := time.Now().Add(time.Hour * 24) //设置过期时间，Add方法用于在当前时间基础上增加指定的时间间隔 当前设置为24小时
	setClaims := &AdminClaims{                   //设置jwt中储存什么的数据
		AdminID:  adminID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()), // 签发时间
			Issuer:    "manager",                      //iss (issuer)：签发人
			//exp (expiration time)：过期时间
			//sub (subject)：主题
			//aud (audience)：受众
			//nbf (Not Before)：生效时间
			//iat (Issued At)：签发时间
			//jti (JWT ID)：编号// 签发者
		},
	}

	// 生成token（使用HS256算法）
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, setClaims)

	// 签名并生成字符串
	tokenString, err := token.SignedString(adminJwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// 解析并校验管理员token
func ParseAdminToken(tokenString string) (*AdminClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("不支持的token签名算法")
		} //token.Method.(*jwt.SigningMethodHMAC) 是 Go 语言的类型断言，
		// 作用是 “判断 token.Method 的实际类型是不是 *jwt.SigningMethodHMAC”。 不需要用到转换后的具体值（value）所以忽略
		return adminJwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	// 校验token有效性
	if claims, ok := token.Claims.(*AdminClaims); ok && token.Valid { //断言，看解析后得到的这个Claim和自定义的AdminClaims是不是一样的，
		return claims, nil
	}
	return nil, errors.New("无效的管理员token")
}
