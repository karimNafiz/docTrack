package upload_sessions

import "time"

type UploadSession struct {
	ID           string    `gorm:"primaryKey"`
	Filename     string    `gorm:"not null"`
	File_size    int64     `gorm:"not null"`
	Chunk_size   int       `gorm:"not null"`
	Total_chunks int       `gorm:"not null"`
	Created_at   time.Time `gorm:"autoCreateTime"`
}
