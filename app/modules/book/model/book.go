package model

import (
	"time"

	"github.com/google/uuid"
)

type BookStatus string

const (
	StatusPending  BookStatus = "pending"
	StatusActive   BookStatus = "active"
	StatusInactive BookStatus = "inactive"
	StatusBanned   BookStatus = "banned"
	StatusDeleted  BookStatus = "deleted"
)

// Book đại diện cho entity sách trong hệ thống
type Book struct {
	ID          uuid.UUID  `json:"id" gorm:"column:id;"`
	CreatedBy   string     `json:"created_by" gorm:"column:created_by;"`
	CreatedAt   time.Time  `json:"created_at" gorm:"column:created_at;"`
	UpdatedBy   string     `json:"updated_by" gorm:"column:updated_by;"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"column:updated_at;"`
	Status      BookStatus `json:"status" gorm:"column:status;"`
	Title       string     `json:"title" gorm:"column:title;"`
	Author      string     `json:"author" gorm:"column:author;"`
	Description string     `json:"description" gorm:"column:description;"`
	Price       float64    `json:"price" gorm:"column:price;"`
	PublishedAt time.Time  `json:"published_at" gorm:"column:published_at;"`
	CoverImage  string     `json:"cover_image" gorm:"column:cover_image;"`
}

// TableName xác định tên bảng trong database
func (Book) TableName() string {
	return "book_books"
}
