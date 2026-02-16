package utils

import "golang.org/x/crypto/bcrypt"


func HashPassword (password string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", nil 
	}

	return string(hash), nil

}

func ComparePassword (hashedPassword string, newPasswordHashed string) error {

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(newPasswordHashed)); err != nil {
		return nil
	}

	return nil

}