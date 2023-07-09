package main

import (
	"log"
	"crypto/rand"
	"github.com/golang-jwt/jwt/v4"
)

// function to generate the random 256 bit for signing the JWT
func generateKey() []byte {
	buffer := make([]byte, 32)	
	_, err := rand.Read(buffer)
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	return buffer
}

// user credentials 
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"` // the hashed password
	Id string `json:"id"` // the user ID
}

// the 'claims' -> JWT payload
type Claims struct {
	Username string `json:"username"`
	UserId string `json:"user_id"`
	jwt.RegisteredClaims // included claims such as exp, ait and jti
}


