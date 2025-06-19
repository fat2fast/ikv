package model

import (
	"time"

	"gorm.io/gorm"
)

// Book đại diện cho entity sách trong hệ thống
type Book struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"size:200;not null" json:"title"`
	Author      string         `gorm:"size:100;not null" json:"author"`
	ISBN        string         `gorm:"size:20;uniqueIndex" json:"isbn"`
	Description string         `gorm:"type:text" json:"description"`
	Price       float64        `gorm:"type:decimal(10,2)" json:"price"`
	PublishedAt time.Time      `json:"published_at"`
	CoverImage  string         `gorm:"size:255" json:"cover_image"`
	CategoryID  uint           `json:"category_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName xác định tên bảng trong database
func (Book) TableName() string {
	return "books"
}
