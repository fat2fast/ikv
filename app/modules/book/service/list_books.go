package bookservice

import (
	"context"

	bookmodel "fat2fast/ikv/modules/book/model"
	"fat2fast/ikv/shared/datatype"
)

// ListBooksQuery đại diện cho query lấy danh sách books
type ListBooksQuery struct {
	Filter *bookmodel.ListBookFilter
}

// IListBooksRepo interface cho repository list operations
type IListBooksRepo interface {
	GetList(ctx context.Context, filter *bookmodel.ListBookFilter) ([]*bookmodel.Book, int64, error)
}

// ListBooksQueryHandler xử lý query lấy danh sách books
type ListBooksQueryHandler struct {
	bookRepo IListBooksRepo
}

// NewListBooksQueryHandler tạo instance mới của ListBooksQueryHandler
func NewListBooksQueryHandler(bookRepo IListBooksRepo) *ListBooksQueryHandler {
	return &ListBooksQueryHandler{bookRepo: bookRepo}
}

// Execute thực thi query lấy danh sách books
func (h *ListBooksQueryHandler) Execute(ctx context.Context, query *ListBooksQuery) (*bookmodel.BookListResponse, error) {
	// Validate và set default values cho filter
	filter := h.normalizeFilter(query.Filter)

	// Lấy danh sách books từ database
	books, total, err := h.bookRepo.GetList(ctx, filter)
	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// Chuyển đổi sang response DTO
	response := bookmodel.ToListResponse(books, total, filter.Page, filter.PerPage)

	return response, nil
}

// normalizeFilter chuẩn hóa filter với default values
func (h *ListBooksQueryHandler) normalizeFilter(filter *bookmodel.ListBookFilter) *bookmodel.ListBookFilter {
	if filter == nil {
		filter = &bookmodel.ListBookFilter{}
	}

	// Set default pagination
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PerPage <= 0 {
		filter.PerPage = 10
	}
	if filter.PerPage > 100 {
		filter.PerPage = 100
	}

	// Set default sorting
	if filter.SortBy == "" {
		filter.SortBy = "created_at"
	}
	if filter.SortOrder == "" {
		filter.SortOrder = "DESC"
	}

	return filter
}
