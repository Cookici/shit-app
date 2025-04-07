package auth

import (
	"errors"
	"record-project/domain/entity"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// 密钥，实际项目中应该从配置中读取
	jwtSecret = []byte("lrh07062003")
	// token有效期，默认7天
	tokenExpireDuration = 7 * 24 * time.Hour
)

// InitJWT 初始化JWT配置
func InitJWT(secret string, expiration time.Duration) {
	jwtSecret = []byte(secret)
	tokenExpireDuration = expiration
}

// Claims 自定义JWT声明
type Claims struct {
	UserID   uint64 `json:"user_id"`
	OpenID   string `json:"open_id"`
	Nickname string `json:"nickname"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(user *entity.User) (string, error) {
	// 创建声明
	claims := Claims{
		UserID:   user.ID,
		OpenID:   user.OpenID,
		Nickname: user.Nickname,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "user_token",
		},
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	return token.SignedString(jwtSecret)
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*Claims, error) {
	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// 验证令牌
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的令牌")
}