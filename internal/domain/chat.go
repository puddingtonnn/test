package domain

import "time"

type Chat struct {
	ID       int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Title    string    `gorm:"not null;size:200" json:"title"`
	CreateAt time.Time `json:"create_at"`

	Messages []Message `gorm:"constraint:OnDelete:CASCADE;" json:"messages,omitempty"`
}

type Message struct {
	ID       int       `gorm:"primaryKey;autoIncrement" json:"id"`
	ChatID   int       `gorm:"not null;index" json:"chat_id"`
	Text     string    `gorm:"not null;size:5000" json:"text"`
	CreateAt time.Time `json:"create_at"`
}
