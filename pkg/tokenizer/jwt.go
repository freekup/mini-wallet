package tokenizer

import (
	"context"
	"errors"

	"github.com/dgrijalva/jwt-go"
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
