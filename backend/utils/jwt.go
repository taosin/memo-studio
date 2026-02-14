package utils

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("memo-studio-secret-key-change-in-production")

func init() {
	// 生产环境必须设置 MEMO_JWT_SECRET
	if v := os.Getenv("MEMO_JWT_SECRET"); v != "" {
		jwtSecret = []byte(v)
	} else if os.Getenv("MEMO_ENV") == "production" {
		log.Fatal("FATAL: 生产环境必须设置 MEMO_JWT_SECRET 环境变量")
	}
}

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT token (默认24小时有效期)
func GenerateToken(userID int, username string, isAdmin bool) (string, error) {
	return GenerateTokenWithExpiry(userID, username, isAdmin, 24*time.Hour)
}

// GenerateTokenWithExpiry 生成带自定义有效期的 JWT token
func GenerateTokenWithExpiry(userID int, username string, isAdmin bool, expiry time.Duration) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析 JWT token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// RefreshToken 刷新 token (延长有效期)
func RefreshToken(tokenString string) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return GenerateTokenWithExpiry(claims.UserID, claims.Username, claims.IsAdmin, 24*time.Hour)
}
