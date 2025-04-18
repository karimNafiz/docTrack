package models

import "time"

type Session struct {
	ID        string    `gorm:"primaryKey;size:64"`
	User_ID   uint      `gorm:"not null;index"`
	Username  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	ExpiresAt time.Time `gorm:"not null"`
}
