package bookrepository

import (
	"context"
	"time"

	bookmodel "fat2fast/ikv/modules/book/model"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Update cập nhật book theo ID
func (r *BookRepository) Update(ctx context.Context, id uuid.UUID, book *bookmodel.Book) error {
	db := r.dbCtx.GetMainConnection()

	// Set updated_at
	book.UpdatedAt = time.Now()

	// Update book
	result := db.WithContext(ctx).Where("id = ?", id).Updates(book)
	if result.Error != nil {
		return errors.WithStack(result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("book not found")
	}

	return nil
}

// UpdateFields cập nhật các fields cụ thể
func (r *BookRepository) UpdateFields(ctx context.Context, id uuid.UUID, fields map[string]interface{}) error {
	db := r.dbCtx.GetMainConnection()

	// Add updated_at
	fields["updated_at"] = time.Now()

	// Update specific fields
	result := db.WithContext(ctx).Model(&bookmodel.Book{}).
		Where("id = ?", id).Updates(fields)

	if result.Error != nil {
		return errors.WithStack(result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("book not found")
	}

	return nil
}

// UpdateStatus cập nhật status của book
func (r *BookRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status bookmodel.BookStatus) error {
	return r.UpdateFields(ctx, id, map[string]interface{}{
		"status": status,
	})
}
