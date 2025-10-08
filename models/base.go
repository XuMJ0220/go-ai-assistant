package models

import "time"

type BaseModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// User 对应 users 表
type User struct {
	BaseModel
	Username     string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	Email        string `gorm:"unique"`
}

// KnowledgeBase 对应 knowledge_bases 表 知识库模型
type KnowledgeBase struct {
	BaseModel
	UserID                 uint   `gorm:"not null"`
	Name                   string `gorm:"not null"`
	Description            string
	VectorDBCollectionName string `gorm:"unique;not null"`
	User                   User   `gorm:"foreignKey:UserID"` // 定义外键关联
}

// Document 对应 documents 表 文档模型
type Document struct {
	BaseModel
	KnowledgeBaseID uint   `gorm:"not null"`
	FileName        string `gorm:"not null"`
	FilePath        string `gorm:"not null"`
	FileSize        int64
	Status          string `gorm:"not null"`
	KnowledgeBase   KnowledgeBase `gorm:"foreignKey:KnowledgeBaseID"` // 定义外键关联
}