package security

import (
	"mini-bank/pkg/utils/apperror"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", apperror.New(apperror.ErrInternalServer, "terjadi kesalahan internal server")
	}

	return string(bytes), nil
}

func ComparePassword(plainPassword string, hashPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(plainPassword))
	return err == nil
}
