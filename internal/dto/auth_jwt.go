package dto

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/kholidss/movie-fest-skilltest/internal/entity"
	"time"
)

func BuildAuthJWTClaims(expiredUnix int64, user *entity.User) jwt.MapClaims {
	return jwt.MapClaims{
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"exp": expiredUnix,
		"user": map[string]any{
			"user_id":     user.ID,
			"user_entity": user.Entity,
		},
	}
}
