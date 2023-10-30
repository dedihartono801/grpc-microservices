package helpers

import (
	"time"

	"github.com/dedihartono801/auth-svc/pkg/config"
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserId uint64 `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

type ClaimsExternal struct {
	Channel      string `json:"channel"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	jwt.StandardClaims
}

func GenerateToken(id uint64, email string) (string, error) {
	expirationTime := time.Now().Add(time.Hour * time.Duration(24))
	claims := &Claims{
		UserId: id,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(config.GetEnv("JWT_SECRET")))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
