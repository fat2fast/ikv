package bookservice

import (
	"context"
	"time"

	bookmodel "fat2fast/ikv/modules/book/model"
	"fat2fast/ikv/shared/datatype"

	"github.com/google/uuid"
)

// CreateBookCommand đại diện cho command tạo book mới
type CreateBookCommand struct {
	Dto bookmodel.CreateBookRequest
}

// ICreateBookRepo interface cho repository create operations
type ICreateBookRepo interface {
	Insert(ctx context.Context, book *bookmodel.Book) error
}

// CreateBookCommandHandler xử lý command tạo book mới
type CreateBookCommandHandler struct {
	bookRepo ICreateBookRepo
}

// NewCreateBookCommandHandler tạo instance mới của CreateBookCommandHandler
func NewCreateBookCommandHandler(bookRepo ICreateBookRepo) *CreateBookCommandHandler {
	return &CreateBookCommandHandler{bookRepo: bookRepo}
}

// Execute thực thi command tạo book mới
func (h *CreateBookCommandHandler) Execute(ctx context.Context, cmd *CreateBookCommand) (*bookmodel.CreateBookResponse, error) {
	// Validate command
	if err := h.validateCreateCommand(cmd); err != nil {
		return nil, err
	}

	// Tạo UUID mới
	newId := uuid.New()
	now := time.Now()

	// Tạo book entity từ command
	book := &bookmodel.Book{
		ID:          newId,
		Title:       cmd.Dto.Title,
		Author:      cmd.Dto.Author,
		Description: cmd.Dto.Description,
		Price:       cmd.Dto.Price,
		PublishedAt: cmd.Dto.PublishedAt,
		CoverImage:  cmd.Dto.CoverImage,
		Status:      bookmodel.StatusActive,
		CreatedBy:   "system", // TODO: Lấy từ context user hiện tại
		CreatedAt:   now,
		UpdatedBy:   "system", // TODO: Lấy từ context user hiện tại
		UpdatedAt:   now,
	}

	// Lưu vào database
	err := h.bookRepo.Insert(ctx, book)
	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// Trả về response
	response := &bookmodel.CreateBookResponse{
		ID: book.ID,
	}

	return response, nil
}

// validateCreateCommand validate command tạo book
func (h *CreateBookCommandHandler) validateCreateCommand(cmd *CreateBookCommand) error {
	// Validate title
	if cmd.Dto.Title == "" {
		return datatype.ErrBadRequest.WithError("Title is required")
	}

	// Validate author
	if cmd.Dto.Author == "" {
		return datatype.ErrBadRequest.WithError("Author is required")
	}

	// Validate price
	if cmd.Dto.Price <= 0 {
		return datatype.ErrBadRequest.WithError("Price must be greater than 0")
	}

	return nil
}
