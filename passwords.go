package main

import (
	"golang.org/x/crypto/bcrypt"
)

// Generate a hash for a given password
func genHash(pswd string) (string, error) {
	data, err_data := bcrypt.GenerateFromPassword([]byte(pswd), 10)
	return string(data), err_data	
}

// Compare the stored hash with the given password
func compHash(hash string, pswd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pswd))
	return err == nil	
}
