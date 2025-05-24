package upload_sessions

import "time"

type UploadSession struct {
	ID             string    `gorm:"primaryKey;column:id"`
	UserID         uint      `gorm:"index;column:user_ID"`
	ParentID       uint      `gorm:"index;column:parent_ID"`
	Status         string    `gorm:"type:text;not null;default:'QUEUED';check:status IN ('QUEUED','IN_PROGRESS','COMPLETED','FAILED','CANCELED')"`
	Filename       string    `gorm:"not null;column:filename"`
	FileSize       int64     `gorm:"not null;column:file_size"`
	ChunkSize      int       `gorm:"not null;column:chunk_size"`
	FileType       string    `gorm:"not null;column:file_type"`
	TotalChunks    int       `gorm:"not null;column:total_chunks"`
	ChunksUploaded int       `gorm:"not null;default 0"`
	CreatedAt      time.Time `gorm:"autoCreateTime;column:created_at"`
}
