package pdfs

type PDF struct {
	ID               uint   `gorm:"primaryKey"`
	UserID           uint   `gorm:"not null;index"`
	OriginalFilename string `gorm:"not null"`
	FileSaveLocation string `gorm:"not null"`
	SizeBytes        int64  `gorm:"not null"` // its better to keep this int64 as in the postgress table its bigint
	// database/sql drivers will give a type mismatch error
	PageNo uint `gorm:"not null"`
}
