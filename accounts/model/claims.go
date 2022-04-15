package model

import "github.com/golang-jwt/jwt"

type UserClaims struct {
	UserId int `json:"uid"`
	jwt.StandardClaims
}
