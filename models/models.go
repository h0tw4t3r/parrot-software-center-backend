package models

import "github.com/dgrijalva/jwt-go"

type User struct {
	ID       int
	Username string
	Password string
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
