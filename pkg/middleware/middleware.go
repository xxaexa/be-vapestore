package middleware

import (
	"clean-architecture/model/dto/json"
	"clean-architecture/model/entity"
	"encoding/base64"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

// basic auth with encode
func BasicAuthEncode(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		json.NewResponseUnauthorized(c, "Authorization header is required", "02", "01")
		c.Abort()
		return
	}

	const prefix = "Basic "
	if !strings.HasPrefix(authHeader, prefix) {
		json.NewResponseUnauthorized(c, "Invalid authorization header", "02", "01")
		c.Abort()
		return
	}

	decoded, err := base64.StdEncoding.DecodeString(authHeader[len(prefix):])
	if err != nil {
		json.NewResponseUnauthorized(c, "Failed to decode authorization header", "02", "01")
		c.Abort()
		return
	}

	creds := strings.SplitN(string(decoded), ":", 2)
	if len(creds) != 2 {
		json.NewResponseUnauthorized(c, "Invalid authorization format", "02", "01")
		c.Abort()
		return
	}

	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("SECRET_CLIENT_ID")

	if clientId != creds[0] || clientSecret != creds[1] {
		json.NewResponseUnauthorized(c, "Invalid credentials", "02", "01")
		c.Abort()
		return
	}

	c.Next()
}

var (
	applicationName  = "incubation-golang"
	jwtSigningMethod = jwt.SigningMethodHS256
	jwtSignatureKey  = []byte(os.Getenv("JWT_SECRET"))
)

// generate token jwt
func GenerateTokenJwt(id string, expiredAt int64) (string, error) {
	loginExpDuration := time.Duration(expiredAt) * time.Minute
	myExpiresAt := time.Now().Add(loginExpDuration).Unix()
	claims := entity.JwtClaim{
		StandardClaims: jwt.StandardClaims{
			Issuer:    applicationName,
			ExpiresAt: myExpiresAt,
		},
		ID: id,
	}

	token := jwt.NewWithClaims(
		jwtSigningMethod,
		claims,
	)

	signedToken, err := token.SignedString(jwtSignatureKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil

}

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			json.NewResponseUnauthorized(c, "Invalid Token", "02", "01")
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", -1))
		claims := &entity.JwtClaim{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSignatureKey, nil
		})
		if err != nil || !token.Valid {
			json.NewResponseUnauthorized(c, "Invalid Token", "02", "01")
			c.Abort()
			return
		}

		c.Set("userID", claims.ID)

		c.Next()
	}
}
