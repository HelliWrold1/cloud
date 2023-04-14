package crypt

import (
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"time"
)

var jwtSecret = []byte(viper.GetString("server.jwtSecret"))

type Claims struct {
	UserId uint `json:"user_id"`
	jwt.StandardClaims
}

// GenerateToken 签发用户Token
func GenerateToken(userID uint) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour) // 过期时间
	claims := Claims{
		UserId: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "38384-SearchEngine",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // 加密
	token, err := tokenClaims.SignedString(jwtSecret)                // 签名
	return token, err
}

// ParseToken 验证用户token
func ParseToken(token string) (*Claims, error) {
	// 解密
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	// 比较
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
