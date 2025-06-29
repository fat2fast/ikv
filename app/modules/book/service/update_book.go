package bookservice

import (
	"context"

	bookmodel "fat2fast/ikv/modules/book/model"
	"fat2fast/ikv/shared/datatype"

	"github.com/google/uuid"
)

// UpdateBookCommand đại diện cho command cập nhật book
type UpdateBookCommand struct {
	ID  uuid.UUID
	Dto bookmodel.UpdateBookRequest
}

// IUpdateBookRepo interface cho repository update operations
type IUpdateBookRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (*bookmodel.Book, error)
	UpdateFields(ctx context.Context, id uuid.UUID, fields map[string]interface{}) error
}

// UpdateBookCommandHandler xử lý command cập nhật book
type UpdateBookCommandHandler struct {
	bookRepo IUpdateBookRepo
}

// NewUpdateBookCommandHandler tạo instance mới của UpdateBookCommandHandler
func NewUpdateBookCommandHandler(bookRepo IUpdateBookRepo) *UpdateBookCommandHandler {
	return &UpdateBookCommandHandler{bookRepo: bookRepo}
}

// Execute thực thi command cập nhật book
func (h *UpdateBookCommandHandler) Execute(ctx context.Context, cmd *UpdateBookCommand) error {
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

	// Prepare update fields
	updateFields := h.buildUpdateFields(&cmd.Dto)

	// Validate update fields
	if err := h.validateUpdateFields(updateFields); err != nil {
		return err
	}

	// Cập nhật book
	err = h.bookRepo.UpdateFields(ctx, cmd.ID, updateFields)
	if err != nil {
		if err.Error() == "book not found" {
			return datatype.ErrNotFound.WithError("Book not found")
		}
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return nil
}

// buildUpdateFields xây dựng map các fields cần update
func (h *UpdateBookCommandHandler) buildUpdateFields(dto *bookmodel.UpdateBookRequest) map[string]interface{} {
	fields := make(map[string]interface{})

	// Set updated_by
	fields["updated_by"] = "system" // TODO: Lấy từ context user hiện tại

	// Chỉ update các fields không empty
	if dto.Title != "" {
		fields["title"] = dto.Title
	}
	if dto.Author != "" {
		fields["author"] = dto.Author
	}
	if dto.Description != "" {
		fields["description"] = dto.Description
	}
	if dto.Price > 0 {
		fields["price"] = dto.Price
	}
	if !dto.PublishedAt.IsZero() {
		fields["published_at"] = dto.PublishedAt
	}
	if dto.CoverImage != "" {
		fields["cover_image"] = dto.CoverImage
	}
	if dto.Status != "" {
		fields["status"] = dto.Status
	}

	return fields
}

// validateUpdateFields validate các fields cần update
func (h *UpdateBookCommandHandler) validateUpdateFields(fields map[string]interface{}) error {
	// Validate price
	if price, exists := fields["price"]; exists {
		if priceFloat, ok := price.(float64); ok && priceFloat <= 0 {
			return datatype.ErrBadRequest.WithError("Price must be greater than 0")
		}
	}

	// Validate status
	if status, exists := fields["status"]; exists {
		if statusStr, ok := status.(string); ok {
			validStatuses := []string{"pending", "active", "inactive", "banned", "deleted"}
			isValid := false
			for _, validStatus := range validStatuses {
				if statusStr == validStatus {
					isValid = true
					break
				}
			}
			if !isValid {
				return datatype.ErrBadRequest.WithError("Invalid status value")
			}
		}
	}

	return nil
}
