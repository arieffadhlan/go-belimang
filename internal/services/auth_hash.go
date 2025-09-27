package services

import "golang.org/x/crypto/bcrypt"

func HashingPassword(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), 7)
	return string(bytes), err
}

func ComparePassword(pass string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass)) == nil
}
