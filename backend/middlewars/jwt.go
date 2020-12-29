package middlewars

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Ismail14098/agyn_test_rest/database/models"
	"github.com/Ismail14098/agyn_test_rest/lib/common"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)


var secretKey []byte

func init() {
	pwd, _ := os.Getwd()
	keyPath := pwd + "/jwtsecret.key"

	key, readErr := ioutil.ReadFile(keyPath)
	if readErr != nil {
		panic("failed to load secret key file")
	}
	secretKey = key
}

func validateToken(tokenString string) (common.JSON, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return secretKey, nil
	})

	if err != nil {
		return common.JSON{}, err
	}

	if !token.Valid {
		return common.JSON{}, errors.New("invalid token")
	}

	return token.Claims.(jwt.MapClaims), nil
}

// JWTMiddleware parses JWT token from cookie and stores data and expires date to the context
// JWT Token can be passed as cookie, or Authorization header
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		// failed to read cookie
		if err != nil {
			// try reading HTTP Header
			authorization := c.Request.Header.Get("Authorization")
			if authorization == "" {
				c.Next()
				return
			}
			sp := strings.Split(authorization, "Bearer ")
			// invalid token
			if len(sp) < 1 {
				c.Next()
				return
			}
			tokenString = sp[1]
		}

		tokenData, err := validateToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		var user models.User
		user.Read(tokenData["user"].(common.JSON))
		c.Set("user", user)
		c.Set("token_expire", tokenData["exp"])
		c.Next()
	}
}
