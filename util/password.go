package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

//使用bcrypt加密密码
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if(err != nil){
		return "", fmt.Errorf("faild to hash password: %w", err)
	}
	return string(hashedPassword), err
}

func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}