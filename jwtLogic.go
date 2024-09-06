package main

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

func (cfg *apiConfig) loadSecret() error {
	godotenv.Load()
	cfg.jwtSecret = os.Getenv("JWT_SECRET")
	return nil
}

func (cfg *apiConfig) createJWT(user User) (string, error) {
	expireTime := time.Now().Add(time.Duration(user.Expires) * time.Second)
	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(expireTime),
		Subject:   string(user.Id),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	sTokenString, err := token.SignedString(cfg.jwtSecret)
	if err != nil {
		log.Printf("Signing failed: %s", err)
		return "", err
	}
	return sTokenString, nil
}

// loadSecret must go into main, after cfg init, add error if no .env is found.
