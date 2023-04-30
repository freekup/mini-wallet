package tokenizer

import (
	"context"
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

const (
	SecretKey = "mini-wallet"
)

// GenerateJWTToken used to generate JWT token base on userXID
func GenerateJWTToken(ctx context.Context, userXID string) (token string, err error) {
	if userXID == "" {
		err = errors.New("userXID is empty")
		return
	}

	tJwt := jwt.New(jwt.SigningMethodHS256)

	claims := tJwt.Claims.(jwt.MapClaims)
	claims["xid"] = userXID

	token, err = tJwt.SignedString([]byte(SecretKey))
	if err != nil {
		return
	}

	return
}

// FetchEchoTokenXID used to get XID body token from Echo Framework
func FetchEchoTokenXID(c echo.Context) (xid string) {
	user := c.Get("user")
	if user == nil {
		return
	}

	claims := user.(*jwt.Token).Claims.(jwt.MapClaims)
	if val, ok := claims["xid"]; ok {
		xid = val.(string)
	}

	return
}
