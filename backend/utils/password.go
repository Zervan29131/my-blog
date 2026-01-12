package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 将明文密码加密
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // 14 是 cost 值，越高越安全但也越慢
	return string(bytes), err
}

// CheckPasswordHash 比对明文密码和数据库的哈希值
// 返回 true 表示密码正确
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
