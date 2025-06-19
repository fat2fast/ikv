package model

import (
	"time"
)

// CreateBookRequest đại diện cho dữ liệu đầu vào khi tạo sách mới
type CreateBookRequest struct {
	Title       string    `json:"title" binding:"required,min=3,max=200"`
	Author      string    `json:"author" binding:"required,min=2,max=100"`
	ISBN        string    `json:"isbn" binding:"required"`
	Description string    `json:"description"`
	Price       float64   `json:"price" binding:"required,min=0.01"`
	PublishedAt time.Time `json:"published_at"`
	CoverImage  string    `json:"cover_image"`
	CategoryID  uint      `json:"category_id"`
}

// UpdateBookRequest đại diện cho dữ liệu đầu vào khi cập nhật sách
type UpdateBookRequest struct {
	Title       string    `json:"title" binding:"omitempty,min=3,max=200"`
	Author      string    `json:"author" binding:"omitempty,min=2,max=100"`
	ISBN        string    `json:"isbn"`
	Description string    `json:"description"`
	Price       float64   `json:"price" binding:"omitempty,min=0.01"`
	PublishedAt time.Time `json:"published_at"`
	CoverImage  string    `json:"cover_image"`
	CategoryID  uint      `json:"category_id"`
}

// BookResponse đại diện cho dữ liệu trả về khi lấy thông tin sách
type BookResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	ISBN        string    `json:"isbn"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	PublishedAt time.Time `json:"published_at"`
	CoverImage  string    `json:"cover_image"`
	CategoryID  uint      `json:"category_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BookListResponse đại diện cho danh sách sách trả về
type BookListResponse struct {
	Books      []BookResponse `json:"books"`
	TotalCount int64          `json:"total_count"`
	Page       int            `json:"page"`
	PerPage    int            `json:"per_page"`
}
