package models

import (
    "github.com/dgrijalva/jwt-go"
    "time"
    config "github.com/jelgar/login/config"
)

// Kinda tempted to make this a wrapper class of what i need from JWT libary so i can just swap it out with my own implemntation later but idk lets see what happens

// Object encoded in JWT to give claim/permissions of user
// jwt.StandardClaims is an embedded type (sounds pretty cool need to google it tbh) but that gives otherfields to the struct such as expiry time
type Claims struct {
    Username string `json:"username"`
    jwt.StandardClaims
}

// Not sure if this is super dumb or super smart but gonna try and roll with it for now
func NewJWT(method jwt.SigningMethod, claims *Claims) *jwt.Token {
    return jwt.NewWithClaims(method, claims)
}

func DefaultSignMethod() *jwt.SigningMethodHMAC {
    return jwt.SigningMethodHS256
}

func NewStandardClaims(time time.Time) jwt.StandardClaims {
    return jwt.StandardClaims{
        ExpiresAt: time.Unix(),
    }
}

func GetKey() []byte {
    return config.SecretKey
}
