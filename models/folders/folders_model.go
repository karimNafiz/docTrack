package folders

import "time"

type Folder struct {
	ID                     uint      `gorm:"primaryKey;column:id"`
	OwnerID                uint      `gorm:"index;column:owner_id"`
	ParentID               uint      `gorm:"index;column:parent_id"`
	Name                   string    `gorm:"not null;column:name"`
	Slug                   string    `gorm:"not null;column:slug"`
	ParentMaterializedPath string    `gorm:"not null;column:parent_materialized_path"`
	Depth                  uint      `gorm:"not null;column:depth"`
	CreatedAt              time.Time `gorm:"autoCreateTime;not null;column:created_at"`
	UpdatedAt              time.Time `gorm:"autoCreateTime;not null;column:updated_at"`
}
