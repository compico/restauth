package middleware

import (
	"github.com/dgrijalva/jwt-go"
)

type JWTConfig struct {
	Data   jwt.MapClaims
	Method jwt.SigningMethod
}

func (config JWTConfig) NewToken() *jwt.Token {
	token := jwt.NewWithClaims(config.Method, config.Data)
	return token
}
