package bookrepository

import (
	"context"

	bookmodel "fat2fast/ikv/modules/book/model"

	"github.com/pkg/errors"
)

// Insert tạo book mới trong database
func (r *BookRepository) Insert(ctx context.Context, book *bookmodel.Book) error {
	db := r.dbCtx.GetMainConnection()

	// Thực hiện insert
	if err := db.WithContext(ctx).Create(book).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
