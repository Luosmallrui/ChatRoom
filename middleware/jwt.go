package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Options jwt.RegisteredClaims

type AuthClaims struct {
	Guard string `json:"guard"` // 授权守卫
	jwt.RegisteredClaims
}

func NewNumericDate(t time.Time) *jwt.NumericDate {
	return jwt.NewNumericDate(t)
}

// GenerateToken 生成 JWT 令牌
func GenerateToken(guard string, secret string, ops *Options) string {

	claims := AuthClaims{
		Guard: guard,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  ops.Audience,
			ExpiresAt: ops.ExpiresAt,
			ID:        ops.ID,
			IssuedAt:  ops.IssuedAt,
			Issuer:    ops.Issuer,
			NotBefore: ops.NotBefore,
			Subject:   ops.Subject,
		},
	}

	tokenString, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))

	return tokenString
}

// ParseToken 解析 JWT Token
func ParseToken(token string, secret string) (*AuthClaims, error) {

	data, err := jwt.ParseWithClaims(token, &AuthClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if claims, ok := data.Claims.(*AuthClaims); ok && data.Valid {
		return claims, nil
	}

	return nil, err
}

const JWTSessionConst = "__JWT_SESSION__"

var (
	ErrNoAuthorize = errors.New("授权异常，请登录后操作! ")
)

type IStorage interface {
	// IsBlackList 判断是否是黑名单
	IsBlackList(ctx context.Context, token string) bool
}

type JSession struct {
	Uid       int    `json:"uid"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

// Auth 授权中间件
func Auth(secret string, guard string, storage IStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := AuthHeaderToken(c)

		claims, err := verify(guard, secret, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": err.Error()})
			return
		}

		if storage.IsBlackList(c.Request.Context(), token) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "请登录再试"})
			return
		}

		uid, err := strconv.Atoi(claims.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "解析 jwt 失败"})
			return
		}

		c.Set(JWTSessionConst, &JSession{
			Uid:       uid,
			Token:     token,
			ExpiresAt: claims.ExpiresAt.Unix(),
		})

		c.Next()
	}
}

func AuthHeaderToken(c *gin.Context) string {
	token := c.GetHeader("Authorization")
	token = strings.TrimSpace(strings.TrimPrefix(token, "Bearer"))

	// Headers 中没有授权信息则读取 url 中的 token
	if token == "" {
		token = c.DefaultQuery("token", "")
	}

	return token
}

func verify(guard string, secret string, token string) (*AuthClaims, error) {

	if token == "" {
		return nil, ErrNoAuthorize
	}

	claims, err := ParseToken(token, secret)
	if err != nil {
		return nil, err
	}

	// 判断权限认证守卫是否一致
	if claims.Guard != guard || claims.Valid() != nil {
		return nil, ErrNoAuthorize
	}

	return claims, nil
}
