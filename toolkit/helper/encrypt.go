package helper

import (
	"golang.org/x/crypto/bcrypt"
)

// Encrypt 给密码加密
func Encrypt(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

// Verify 校验密码
func Verify(hashedPwd string, inputPwd string) bool {
	// Returns true on success, hashedPwd is for the database.
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(inputPwd))
	if err != nil {
		return false
	} else {
		return true
	}
}
