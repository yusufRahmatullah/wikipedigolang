package util

import (
	"golang.org/x/crypto/bcrypt"
)

// HashAndSalt generate Hash and salt for password
func HashAndSalt(pwd []byte) (string, error) {
	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)

	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), err
}

// ComparePasswords compare hashed password with plain password
func ComparePasswords(hashedPwd string, plainPwd []byte) error {

	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)

	return bcrypt.CompareHashAndPassword(byteHash, plainPwd)
}
