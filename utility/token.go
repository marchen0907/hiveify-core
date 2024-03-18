package utility

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/golang-jwt/jwt/v5"
	"hiveify-core/internal/model/entity"
	"hiveify-core/internal/response"
	"time"
)

// CustomClaims 自定义JWT声明
type CustomClaims struct {
	*entity.User
	jwt.RegisteredClaims
}

const (
	secretKey = "Yw#Lwpak#525178"
)

var ctx = context.Background()

// GetToken 生成token
func GetToken(user *entity.User) (string, error) {
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
		g.Log().Errorf(ctx, "Error generating token: %s", err.Error())
		return "", gerror.NewCodef(response.CodeServerError, "生成token失败!")
	}
	return tokenString, nil
}

// ParseToken 解析token
func ParseToken(ctx context.Context, tokenString string) (*entity.User, error) {
	// 解析JWT
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证方法，检查签名方法是否正确
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			g.Log().Warningf(ctx, "Unexpected signing method: %v", token.Header["alg"])
			return nil, gerror.NewCodef(response.CodeNotAuthorized, "token非法！")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		g.Log().Errorf(ctx, "Error parsing token: %s", err.Error())
		return nil, gerror.NewCodef(response.CodeServerError, "token解析失败!")
	}

	// 验证token并检查是否解析成功
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims.User, nil
	} else {
		g.Log().Errorf(ctx, "Error parsing token: %s", err)
		return nil, gerror.NewCodef(response.CodeNotAuthorized, "token非法")
	}
}

// ValidateToken 验证token是否有效
func ValidateToken(ctx context.Context, tokenString string) bool {
	_, err := ParseToken(ctx, tokenString)
	if err != nil {
		return false
	}
	return true
}
