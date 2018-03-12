package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TODO: this is utterly insecure, get this key securely from our infra tools
var sampleSecret = []byte("insecure")

var method = jwt.SigningMethodHS256

type TokenClaims struct {
	User string `json:"user"`
	jwt.StandardClaims
}

func TokenForUser(user string) (string,error) {
	// we currently allow tokens to be used forever but embed issue time so we can change that later
	claims := TokenClaims{User: user, StandardClaims: jwt.StandardClaims{ IssuedAt: time.Now().Unix()}}
	return jwt.NewWithClaims(method, claims).SignedString(sampleSecret)
}

func JwtKeyFunc(token *jwt.Token) (interface{}, error) {
		if token.Method != method {
			return nil, fmt.Errorf("invalid alg") // require validation with the same alg we used to generate signature
		}
		return sampleSecret, nil
}