package services

import (
	"errors"
	"my-blog-backend/database"
	"my-blog-backend/models"
	"my-blog-backend/utils"

	"golang.org/x/crypto/bcrypt"
)

// Login 登录逻辑
func Login(username, password string) (string, error) {
	var user models.User

	// 1. 查找用户
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return "", errors.New("用户不存在")
	}

	// 2. 验证密码
	// 注意：注册时必须使用 bcrypt.GenerateFromPassword 加密密码存入数据库
	// 如果你数据库里现在存的是明文 "123456"，这里比对会失败。
	// 临时测试如果是明文，可以改用: if user.Password != password { ... }
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("密码错误")
	}

	// 3. 生成 Token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}
