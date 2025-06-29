package model

import (
	"time"

	"github.com/google/uuid"
)

// CreateBookRequest đại diện cho dữ liệu đầu vào khi tạo sách mới
type CreateBookRequest struct {
	Title       string    `json:"title" binding:"required,min=3,max=200"`
	Author      string    `json:"author" binding:"required,min=2,max=100"`
	Description string    `json:"description" binding:"max=1000"`
	Price       float64   `json:"price" binding:"required,min=0.01"`
	PublishedAt time.Time `json:"published_at"`
	CoverImage  string    `json:"cover_image" binding:"omitempty,url"`
}

// UpdateBookRequest đại diện cho dữ liệu đầu vào khi cập nhật sách
type UpdateBookRequest struct {
	Title       string    `json:"title" binding:"omitempty,min=3,max=200"`
	Author      string    `json:"author" binding:"omitempty,min=2,max=100"`
	Description string    `json:"description" binding:"omitempty,max=1000"`
	Price       float64   `json:"price" binding:"omitempty,min=0.01"`
	PublishedAt time.Time `json:"published_at"`
	CoverImage  string    `json:"cover_image" binding:"omitempty,url"`
	Status      string    `json:"status" binding:"omitempty,oneof=pending active inactive banned deleted"`
}

// BookResponse đại diện cho dữ liệu trả về khi lấy thông tin sách
type BookResponse struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Author      string     `json:"author"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	PublishedAt time.Time  `json:"published_at"`
	CoverImage  string     `json:"cover_image"`
	Status      BookStatus `json:"status"`
	CreatedBy   string     `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedBy   string     `json:"updated_by"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// BookListResponse đại diện cho dữ liệu trả về khi lấy danh sách sách
type BookListResponse struct {
	Items      []*BookResponse `json:"items"`
	TotalCount int64           `json:"total_count"`
	Page       int             `json:"page"`
	PerPage    int             `json:"per_page"`
	TotalPages int             `json:"total_pages"`
}

// ListBookFilter đại diện cho bộ lọc khi lấy danh sách sách
type ListBookFilter struct {
	Page        int       `json:"page" form:"page" binding:"omitempty,min=1"`
	PerPage     int       `json:"per_page" form:"per_page" binding:"omitempty,min=1,max=100"`
	Status      string    `json:"status" form:"status" binding:"omitempty,oneof=pending active inactive banned deleted"`
	Search      string    `json:"search" form:"search" binding:"omitempty,max=100"`
	Author      string    `json:"author" form:"author" binding:"omitempty,max=100"`
	SortBy      string    `json:"sort_by" form:"sort_by" binding:"omitempty,oneof=title author price created_at updated_at"`
	SortOrder   string    `json:"sort_order" form:"sort_order" binding:"omitempty,oneof=ASC DESC"`
	CreatedFrom time.Time `json:"created_from" form:"created_from"`
	CreatedTo   time.Time `json:"created_to" form:"created_to"`
	PriceMin    float64   `json:"price_min" form:"price_min" binding:"omitempty,min=0"`
	PriceMax    float64   `json:"price_max" form:"price_max" binding:"omitempty,min=0"`
}

// CreateBookResponse đại diện cho dữ liệu trả về khi tạo sách mới
type CreateBookResponse struct {
	ID uuid.UUID `json:"id"`
}

// Conversion methods

// ToResponse chuyển đổi Book entity sang BookResponse
func (b *Book) ToResponse() *BookResponse {
	return &BookResponse{
		ID:          b.ID,
		Title:       b.Title,
		Author:      b.Author,
		Description: b.Description,
		Price:       b.Price,
		PublishedAt: b.PublishedAt,
		CoverImage:  b.CoverImage,
		Status:      b.Status,
		CreatedBy:   b.CreatedBy,
		CreatedAt:   b.CreatedAt,
		UpdatedBy:   b.UpdatedBy,
		UpdatedAt:   b.UpdatedAt,
	}
}

// ToListResponse chuyển đổi danh sách Book entities sang BookListResponse
func ToListResponse(books []*Book, total int64, page, perPage int) *BookListResponse {
	items := make([]*BookResponse, len(books))
	for i, book := range books {
		items[i] = book.ToResponse()
	}

	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	return &BookListResponse{
		Items:      items,
		TotalCount: total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}
}
