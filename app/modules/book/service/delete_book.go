package bookservice

import (
	"context"

	bookmodel "fat2fast/ikv/modules/book/model"
	"fat2fast/ikv/shared/datatype"

	"github.com/google/uuid"
)

// DeleteBookCommand đại diện cho command xóa book
type DeleteBookCommand struct {
	ID   uuid.UUID
	Soft bool // true = soft delete, false = hard delete
}

// IDeleteBookRepo interface cho repository delete operations
type IDeleteBookRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (*bookmodel.Book, error)
	Delete(ctx context.Context, id uuid.UUID) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

// DeleteBookCommandHandler xử lý command xóa book
type DeleteBookCommandHandler struct {
	bookRepo IDeleteBookRepo
}

// NewDeleteBookCommandHandler tạo instance mới của DeleteBookCommandHandler
func NewDeleteBookCommandHandler(bookRepo IDeleteBookRepo) *DeleteBookCommandHandler {
	return &DeleteBookCommandHandler{bookRepo: bookRepo}
}

// Execute thực thi command xóa book
func (h *DeleteBookCommandHandler) Execute(ctx context.Context, cmd *DeleteBookCommand) error {
	// Validate command
	if cmd.ID == uuid.Nil {
		return datatype.ErrBadRequest.WithError("Book ID is required")
	}

	// Kiểm tra book có tồn tại không
	_, err := h.bookRepo.GetByID(ctx, cmd.ID)
	if err != nil {
		if err.Error() == "book not found" {
			return datatype.ErrNotFound.WithError("Book not found")
		}
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// Thực hiện xóa book
	if cmd.Soft {
		err = h.bookRepo.SoftDelete(ctx, cmd.ID)
	} else {
		err = h.bookRepo.Delete(ctx, cmd.ID)
	}

	if err != nil {
		if err.Error() == "book not found" {
			return datatype.ErrNotFound.WithError("Book not found")
		}
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return nil
}
