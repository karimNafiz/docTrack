package upload_sessions

import "time"

type UploadSession struct {
	ID          string    `gorm:"primaryKey;column:id"`
	UserID      uint      `gorm:"index;column:user_ID"`
	ParentID    uint      `gorm:"index;column:parent_ID"`
	Filename    string    `gorm:"not null;column:filename"`
	FileSize    int64     `gorm:"not null;column:file_size"`
	ChunkSize   int       `gorm:"not null;column:chunk_size"`
	FileType    string    `gorm:"not null;column:file_type"`
	TotalChunks int       `gorm:"not null;column:total_chunks"`
	CreatedAt   time.Time `gorm:"autoCreateTime;column:created_at"`
}
