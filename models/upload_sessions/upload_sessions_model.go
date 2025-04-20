package upload_sessions

import "time"

type UploadSession struct {
	ID          string    `gorm:"primaryKey;column:id"`
	Filename    string    `gorm:"not null;column:filename"`
	FileSize    int64     `gorm:"not null;column:file_size"`
	ChunkSize   int       `gorm:"not null;column:chunk_size"`
	TotalChunks int       `gorm:"not null;column:total_chunks"`
	CreatedAt   time.Time `gorm:"autoCreateTime;column:created_at"`
}
