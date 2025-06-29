package bookrepository

import (
	"context"
	"time"

	bookmodel "fat2fast/ikv/modules/book/model"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Delete xóa vĩnh viễn book
func (r *BookRepository) Delete(ctx context.Context, id uuid.UUID) error {
	db := r.dbCtx.GetMainConnection()

	result := db.WithContext(ctx).Where("id = ?", id).Delete(&bookmodel.Book{})
	if result.Error != nil {
		return errors.WithStack(result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("book not found")
	}

	return nil
}

// SoftDelete xóa mềm book (set status = deleted)
func (r *BookRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	db := r.dbCtx.GetMainConnection()

	now := time.Now()
	result := db.WithContext(ctx).Model(&bookmodel.Book{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     bookmodel.StatusDeleted,
			"updated_at": now,
		})

	if result.Error != nil {
		return errors.WithStack(result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("book not found")
	}

	return nil
}
