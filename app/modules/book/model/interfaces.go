package model

import (
	"context"

	"github.com/google/uuid"
)

// Repository interfaces cho Book

// ICreateBookRepository interface cho create operations
type ICreateBookRepository interface {
	Insert(ctx context.Context, book *Book) error
}

// IReadBookRepository interface cho read operations
type IReadBookRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Book, error)
	GetList(ctx context.Context, filter *ListBookFilter) ([]*Book, int64, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}

// IUpdateBookRepository interface cho update operations
type IUpdateBookRepository interface {
	Update(ctx context.Context, id uuid.UUID, book *Book) error
	UpdateFields(ctx context.Context, id uuid.UUID, fields map[string]interface{}) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status BookStatus) error
}

// IDeleteBookRepository interface cho delete operations
type IDeleteBookRepository interface {
	Delete(ctx context.Context, id uuid.UUID) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

// IBookRepository composite interface cho tất cả CRUD operations
type IBookRepository interface {
	ICreateBookRepository
	IReadBookRepository
	IUpdateBookRepository
	IDeleteBookRepository
}
