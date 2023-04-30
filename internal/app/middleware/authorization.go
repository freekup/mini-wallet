package middleware

import (
	"github.com/freekup/mini-wallet/pkg/tokenizer"
	"github.com/labstack/echo/v4/middleware"
)

var JWTAuth = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningKey: []byte(tokenizer.SecretKey),
	AuthScheme: "Token",
})
