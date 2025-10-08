// models/chat.go
package models

// ChatSession 对应 chat_sessions 表
type ChatSession struct {
	BaseModel
	KnowledgeBaseID uint          `gorm:"not null"`
	UserID          uint          `gorm:"not null"`
	SessionName     string        `gorm:"not null"`
	User            User          `gorm:"foreignKey:UserID"`
	KnowledgeBase   KnowledgeBase `gorm:"foreignKey:KnowledgeBaseID"`
}

// ChatMessage 对应 chat_messages 表
type ChatMessage struct {
	BaseModel
	SessionID uint        `gorm:"not null"`
	Role      string      `gorm:"not null"` // "user" or "assistant"
	Content   string      `gorm:"type:text;not null"`
	Session   ChatSession `gorm:"foreignKey:SessionID"`
}
