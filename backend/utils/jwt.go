package utils // 建议放在 utils 包下

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5" // 注意这里引入的是 v5 版本
)

const TokenExpireDuration = time.Hour * 2

// MySecret 用于签名的密钥，生产环境建议从配置文件或环境变量读取
var MySecret = []byte("5408")

// MyClaims 自定义声明结构体并内嵌 jwt.RegisteredClaims
// jwt-go v4/v5 版本中，StandardClaims 已被 RegisteredClaims 替代
type MyClaims struct {
	UserID   uint   `json:"userid"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenToken 生成 JWT
func GenToken(userid uint, username string) (string, error) {
	// 创建声明
	c := MyClaims{
		UserID:   userid,
		Username: username, // ✅ 修正：使用传入的变量
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)), // 过期时间
			Issuer:    "midNight",                                              // 签发人
		},
	}

	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	// 使用指定的 secret 签名并获得完整的编码后的字符串 token
	return token.SignedString(MySecret)
}

// ParseToken 解析 JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析 token
	// ✅ 修正：使用 ParseWithClaims 解析到具体的结构体
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		// ✅ 安全检查：确保签名方法是我们预期的 HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return MySecret, nil
	})

	if err != nil {
		return nil, err
	}

	// 校验 token 并提取 claims
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
