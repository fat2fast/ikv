package bookrepository

import (
	bookmodel "fat2fast/ikv/modules/book/model"
	sharedinfras "fat2fast/ikv/shared/infras"
)

// BookRepository chứa các phương thức truy cập dữ liệu cho Book
type BookRepository struct {
	dbCtx sharedinfras.IDbContext
}

// NewBookRepository tạo instance mới của BookRepository
func NewBookRepository(dbCtx sharedinfras.IDbContext) bookmodel.IBookRepository {
	return &BookRepository{dbCtx: dbCtx}
}

// GetDBContext trả về database context
func (r *BookRepository) GetDBContext() sharedinfras.IDbContext {
	return r.dbCtx
}
