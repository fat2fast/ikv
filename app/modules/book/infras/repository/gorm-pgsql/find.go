package bookrepository

import (
	"context"
	"strings"

	bookmodel "fat2fast/ikv/modules/book/model"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// GetByID lấy book theo ID
func (r *BookRepository) GetByID(ctx context.Context, id uuid.UUID) (*bookmodel.Book, error) {
	db := r.dbCtx.GetMainConnection()
	var book bookmodel.Book

	err := db.WithContext(ctx).Where("id = ?", id).First(&book).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("book not found")
		}
		return nil, errors.WithStack(err)
	}

	return &book, nil
}

// GetList lấy danh sách books với filter và pagination
func (r *BookRepository) GetList(ctx context.Context, filter *bookmodel.ListBookFilter) ([]*bookmodel.Book, int64, error) {
	db := r.dbCtx.GetMainConnection()
	var books []*bookmodel.Book
	var total int64

	// Build base query
	query := db.WithContext(ctx).Model(&bookmodel.Book{})

	// Apply filters
	query = r.applyFilters(query, filter)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}

	// Apply pagination and sorting
	query = r.applyPaginationAndSorting(query, filter)

	// Execute query
	if err := query.Find(&books).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}

	return books, total, nil
}

// Exists kiểm tra book có tồn tại không
func (r *BookRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	db := r.dbCtx.GetMainConnection()
	var count int64

	err := db.WithContext(ctx).Model(&bookmodel.Book{}).
		Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, errors.WithStack(err)
	}

	return count > 0, nil
}

// applyFilters áp dụng các filter vào query
func (r *BookRepository) applyFilters(query *gorm.DB, filter *bookmodel.ListBookFilter) *gorm.DB {
	// Filter by status
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// Search by title or author
	if filter.Search != "" {
		searchTerm := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where("(LOWER(title) LIKE ? OR LOWER(author) LIKE ?)", searchTerm, searchTerm)
	}

	// Filter by author
	if filter.Author != "" {
		query = query.Where("LOWER(author) LIKE ?", "%"+strings.ToLower(filter.Author)+"%")
	}

	// Filter by date range
	if !filter.CreatedFrom.IsZero() {
		query = query.Where("created_at >= ?", filter.CreatedFrom)
	}
	if !filter.CreatedTo.IsZero() {
		query = query.Where("created_at <= ?", filter.CreatedTo)
	}

	// Filter by price range
	if filter.PriceMin > 0 {
		query = query.Where("price >= ?", filter.PriceMin)
	}
	if filter.PriceMax > 0 {
		query = query.Where("price <= ?", filter.PriceMax)
	}

	return query
}

// applyPaginationAndSorting áp dụng pagination và sorting
func (r *BookRepository) applyPaginationAndSorting(query *gorm.DB, filter *bookmodel.ListBookFilter) *gorm.DB {
	// Sorting
	sortBy := filter.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}

	sortOrder := filter.SortOrder
	if sortOrder == "" {
		sortOrder = "DESC"
	}

	query = query.Order(sortBy + " " + sortOrder)

	// Pagination
	if filter.Page > 0 && filter.PerPage > 0 {
		offset := (filter.Page - 1) * filter.PerPage
		query = query.Offset(offset).Limit(filter.PerPage)
	}

	return query
}
