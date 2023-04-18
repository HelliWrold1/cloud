package crypt

import (
	"golang.org/x/crypto/bcrypt"
)

func GenerateSaltPwd(password string) (string, error) {
	// 加盐处理
	_, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPwd), nil
}

func CheckSaltPwd(password string, hashedPwd string) error {
	// 判断加盐处理后的密码是否正确
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
