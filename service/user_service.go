package service

import (
	"errors"
	"go-ai-assistant/core"
	"go-ai-assistant/models"

	"golang.org/x/crypto/bcrypt"
)

// RegisterUserInput 定义了注册时需要传入的参数
type RegisterUserInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// Register 处理用户注册的业务逻辑
func Register(input RegisterUserInput) (models.User, error) {
	// 检查用户名是否存在
	var existingUser models.User
	// 检查到用户已经存在
	if err := core.DB.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
		return models.User{}, errors.New("username already exists")
	}
	// 用bcrypt处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}
	// 创建用户实例
	newUser := models.User{
		Username:     input.Username,
		PasswordHash: string(hashedPassword),
		Email:        input.Email,
	}
	// 存入数据库
	if err := core.DB.Create(&newUser).Error; err != nil {
		return models.User{}, err
	}
	// 返回用户实例
	return newUser, nil
}
