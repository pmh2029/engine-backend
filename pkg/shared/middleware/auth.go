package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"engine/pkg/shared/auth"
)

func CheckAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			c.JSON(http.StatusUnauthorized, "missing")
			c.Abort()
			return
		}
		tokenReq := strings.Replace(authorization, "Bearer ", "", -1)
		tokenReqParts := strings.Split(tokenReq, ".")
		if len(tokenReqParts) != 3 {
			c.JSON(http.StatusUnauthorized,
				gin.H{"Message": "Token is not JWT"})
			c.Abort()
			return
		}
		if !auth.VerifyJWT(tokenReq) {
			c.JSON(http.StatusUnauthorized,
				gin.H{"Message": "Token is invalid"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func GetUserInfoFromToken(c *gin.Context) (string, error) {
	authorization := c.Request.Header.Get("Authorization")
	if authorization == "" {
		return "", errors.New("unauthorized")
	}

	token := strings.Replace(authorization, "Bearer ", "", -1)
	if !auth.VerifyJWT(token) {
		return "", errors.New("token is not valid")
	}

	decodedToken, err := auth.DecodeJWT(token)
	if err != nil {
		return "", errors.New("error while decoding token")
	}

	email := decodedToken.Claims.(jwt.MapClaims)["email"].(string)

	return email, nil
}
