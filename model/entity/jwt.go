package entity

import "github.com/dgrijalva/jwt-go"

type (
	JwtClaim struct {
		jwt.StandardClaims
		ID    string `json:"id"`
		Roles string `json:"role"`
	}
)
