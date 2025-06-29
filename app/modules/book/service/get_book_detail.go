package bookservice

import (
	"context"

	bookmodel "fat2fast/ikv/modules/book/model"
	"fat2fast/ikv/shared/datatype"

	"github.com/google/uuid"
)

// GetBookDetailQuery đại diện cho query lấy chi tiết book
type GetBookDetailQuery struct {
	ID uuid.UUID
}

// IGetBookDetailRepo interface cho repository read operations
type IGetBookDetailRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (*bookmodel.Book, error)
}

// GetBookDetailQueryHandler xử lý query lấy chi tiết book
type GetBookDetailQueryHandler struct {
	bookRepo IGetBookDetailRepo
}

// NewGetBookDetailQueryHandler tạo instance mới của GetBookDetailQueryHandler
func NewGetBookDetailQueryHandler(bookRepo IGetBookDetailRepo) *GetBookDetailQueryHandler {
	return &GetBookDetailQueryHandler{bookRepo: bookRepo}
}

// Execute thực thi query lấy chi tiết book
func (h *GetBookDetailQueryHandler) Execute(ctx context.Context, query *GetBookDetailQuery) (*bookmodel.BookResponse, error) {
	// Validate query
	if query.ID == uuid.Nil {
		return nil, datatype.ErrBadRequest.WithError("Book ID is required")
	}

	// Lấy book từ database
	book, err := h.bookRepo.GetByID(ctx, query.ID)
	if err != nil {
		if err.Error() == "book not found" {
			return nil, datatype.ErrNotFound.WithError("Book not found")
		}
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// Chuyển đổi sang response DTO
	response := book.ToResponse()

	return response, nil
}
