package pdfs

type PDF struct {
	ID                 uint   `gorm:"primaryKey"`
	User_ID            uint   `gorm:"not null;index"`
	Original_Filename  string `gorm:"not null"`
	File_save_location string `gorm:"not null"`
	Size_Bytes         int64  `gorm:"not null"` // its better to keep this int64 as in the postgress table its bigint
	// database/sql drivers will give a type mismatch error
	Page_No uint `gorm:"not null"`
}
