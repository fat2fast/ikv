package usermodel

import (
	"time"

	"github.com/google/uuid"
)

type LoginForm struct {
	Username string `json:"username" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterForm struct {
	ID        string `json:"-"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8,max=32"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
}

type RegisterResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     string    `json:"phone"`
}

// ProfileResponse đại diện cho dữ liệu profile trả về
type ProfileResponse struct {
	ID        uuid.UUID  `json:"id"`
	Email     string     `json:"email"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Phone     string     `json:"phone"`
	FullName  string     `json:"full_name"`
	Role      UserRole   `json:"role"`
	Status    UserStatus `json:"status"`
	Type      UserType   `json:"type"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// UpdateProfileRequest đại diện cho dữ liệu cập nhật profile
type UpdateProfileRequest struct {
	FirstName string `json:"first_name" binding:"omitempty,min=1,max=50"`
	LastName  string `json:"last_name" binding:"omitempty,min=1,max=50"`
	Phone     string `json:"phone" binding:"omitempty,min=10,max=15"`
}

// CreateBookRequest đại diện cho dữ liệu đầu vào khi tạo sách mới
type CreateUserRequest struct {
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
type UpdateUserRequest struct {
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
