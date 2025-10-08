package service

import (
	"errors"
	"go-ai-assistant/config"
	"go-ai-assistant/core"
	"go-ai-assistant/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUserInput 定义了注册时需要传入的参数
type RegisterUserInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// LoginUserInput 定义了登录时需要传入的参数
type LoginUserInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
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

// Login 处理用户登录的业务逻辑
func Login(input LoginUserInput) (string, error) {
	// 1. 根据用户名查找用户
	var user models.User
	if err := core.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		return "", errors.New("invalid username or password")
	}

	// 2. 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		// 如果密码不匹配，为了安全，我们仍然返回与用户名不存在时相同的错误信息
		return "", errors.New("invalid username or password")
	}

	// 3. 密码验证成功，生成 JWT Token
	token, err := generateJWT(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

// generateJWT 生成 JWT Token
func generateJWT(user models.User) (string, error) {
	// 创建一个声明 (Claims)
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * time.Duration(config.AppConfig.JWT.ExpireHours)).Unix(), // 过期时间
		"iat":      time.Now().Unix(),                                                                  // 签发时间
	}

	// 使用 HS256 签名算法创建一个新的 Token 对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用我们在配置文件中设置的密钥来为 Token 签名，并获取完整的 Token 字符串
	signedToken, err := token.SignedString([]byte(config.AppConfig.JWT.SecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
