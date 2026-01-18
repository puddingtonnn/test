package domain

import "time"

type Chat struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     string    `gorm:"not null;size:200" json:"title"`
	CreatedAt time.Time `json:"created_at"`

	Messages []Message `gorm:"constraint:OnDelete:CASCADE;" json:"messages,omitempty"`
}

type Message struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	ChatID    uint      `gorm:"not null;index" json:"chat_id"`
	Text      string    `gorm:"not null;size:5000" json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
