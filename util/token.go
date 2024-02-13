package util

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"hiveify-core/database"
	"time"
)

// CustomClaims 自定义JWT声明
type CustomClaims struct {
	database.User
	jwt.RegisteredClaims
}

const (
	secretKey = "Yw#Lwpak#525178"
)

// GetToken 生成token
func GetToken(user database.User) (string, error) {
	customClaims := CustomClaims{
		user,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Subject:   "hiveify-core-core",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, customClaims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Errorf("Error generating token: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

// ParseToken 解析token
func parseToken(tokenString string) (database.User, error) {
	// 解析JWT
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证方法，检查签名方法是否正确
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return database.User{}, fmt.Errorf("failed to parse token: %w", err)
	}

	// 验证token并检查是否解析成功
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims.User, nil
	} else {
		return database.User{}, fmt.Errorf("token is not valid: %w", err)
	}
}

// ValidateToken 验证token是否有效
func ValidateToken(tokenString string) bool {
	_, err := parseToken(tokenString)
	if err != nil {
		return false
	}
	return true
}

// GetUserByString 从token字符串中获取用户信息
func GetUserByString(tokenString string) (database.User, error) {
	return parseToken(tokenString)
}

// GetUserByRequest 从HTTP请求中获取用户信息
func GetUserByRequest(c *gin.Context) (database.User, error) {
	token, _ := c.Cookie("Authorization")
	return parseToken(token)
}
