package caddyshardrouter

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var localKeyBytes, _ = os.ReadFile("cert/id_rsa.pub")

func ParseJWT(tokenStr string) (claims jwt.MapClaims, err error) {
	pub, err := jwt.ParseRSAPublicKeyFromPEM(localKeyBytes)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return pub, nil
	})
	if err != nil {
		return nil, err
	}

	var ok bool
	if claims, ok = token.Claims.(jwt.MapClaims); !ok {
		return nil, fmt.Errorf("failed to type assert the jwt claims")
	}
	return
}
