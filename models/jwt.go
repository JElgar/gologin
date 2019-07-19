package models

import (
    "github.com/dgrijalva/jwt-go"
    "time"
    config "github.com/jelgar/login/config"
    "github.com/gin-gonic/gin"
)

// Kinda tempted to make this a wrapper class of what i need from JWT libary so i can just swap it out with my own implemntation later but idk lets see what happens

// Object encoded in JWT to give claim/permissions of user
// jwt.StandardClaims is an embedded type (sounds pretty cool need to google it tbh) but that gives otherfields to the struct such as expiry time

var jwtKey []byte = config.SecretKey

var ErrSignatureInvalid error = jwt.ErrSignatureInvalid

type Claims struct {
    Username string `json:"username"`
    jwt.StandardClaims
}

// Not sure if this is super dumb or super smart but gonna try and roll with it for now

// Create a new JWT with claims, and signing method (default is below)
func NewJWT(method jwt.SigningMethod, claims *Claims) *jwt.Token {
    return jwt.NewWithClaims(method, claims)
}

// Return the default method used to sign JWT in this application
func DefaultSignMethod() *jwt.SigningMethodHMAC {
    return jwt.SigningMethodHS256
}

func NewStandardClaims(time time.Time) jwt.StandardClaims {
    return jwt.StandardClaims{
        ExpiresAt: time.Unix(),
    }
}

// Retruns the JWT secret key used for signing the token
func GetKey() []byte {
    return jwtKey
}

// Checks that the token is valid based on the key
func ParseWClaims(tokenString string, claims jwt.Claims) (*jwt.Token, error) {
    return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
}

//func Refresh

// Need to write a validate function so i dont have to redo this every time
//func ValidateToken(g gin.HandlerFunc) gin.HandlerFunc {
//    return g
//}
