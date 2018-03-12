package server

import (
	"fmt"

	"github.com/bayesianmind/demo-file-server/auth"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const userKey = "user"

func requireValidSession(c *gin.Context) {
	tokenEnc := c.GetHeader("X-Session")
	claims := &auth.TokenClaims{}
	token, err := jwt.ParseWithClaims(tokenEnc, claims, auth.JwtKeyFunc)
	if err != nil {
		fmt.Println("Invalid token: ", err)
	}
	if err != nil || !token.Valid {
		c.AbortWithStatus(403)
		c.String(403, "Invalid session token") // we don't pass back the underlying error as it could be a security risk
	}
	c.Set(userKey, claims.User)
}
