package util

import "golang.org/x/crypto/bcrypt"

// hashing password
func HashingPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// comparing password
func ComparePassword(oldPassword, newPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(oldPassword), []byte(newPassword))
}
