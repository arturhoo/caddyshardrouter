package caddyshardrouter

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var localKeyBytes, _ = os.ReadFile("cert/id_rsa.pub")

func ParseJWT(tokenStr string) (claims jwt.MapClaims) {
	block, _ := pem.Decode(localKeyBytes)
	pub, _ := x509.ParsePKCS1PublicKey(block.Bytes)

	token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return &pub, nil
	})

	claims, _ = token.Claims.(jwt.MapClaims)
	return
}
