package bookhttpgin

import (
	"context"

	bookmodel "fat2fast/ikv/modules/book/model"
	bookservice "fat2fast/ikv/modules/book/service"
)

// Interface definitions cho command handlers
type ICreateBookCommandHandler interface {
	Execute(ctx context.Context, cmd *bookservice.CreateBookCommand) (*bookmodel.CreateBookResponse, error)
}

type IUpdateBookCommandHandler interface {
	Execute(ctx context.Context, cmd *bookservice.UpdateBookCommand) error
}

type IDeleteBookCommandHandler interface {
	Execute(ctx context.Context, cmd *bookservice.DeleteBookCommand) error
}

// Interface definitions cho query handlers
type IGetBookDetailQueryHandler interface {
	Execute(ctx context.Context, query *bookservice.GetBookDetailQuery) (*bookmodel.BookResponse, error)
}

type IListBooksQueryHandler interface {
	Execute(ctx context.Context, query *bookservice.ListBooksQuery) (*bookmodel.BookListResponse, error)
}

// BookHTTPController chứa tất cả handlers cho book CRUD operations
type BookHTTPController struct {
	// Command handlers
	createCmdHdl ICreateBookCommandHandler
	updateCmdHdl IUpdateBookCommandHandler
	deleteCmdHdl IDeleteBookCommandHandler

	// Query handlers
	getDetailQryHdl IGetBookDetailQueryHandler
	listQryHdl      IListBooksQueryHandler
}

// NewBookHTTPController tạo instance mới của BookHTTPController
func NewBookHTTPController(
	createCmdHdl ICreateBookCommandHandler,
	updateCmdHdl IUpdateBookCommandHandler,
	deleteCmdHdl IDeleteBookCommandHandler,
	getDetailQryHdl IGetBookDetailQueryHandler,
	listQryHdl IListBooksQueryHandler,
) *BookHTTPController {
	return &BookHTTPController{
		createCmdHdl:    createCmdHdl,
		updateCmdHdl:    updateCmdHdl,
		deleteCmdHdl:    deleteCmdHdl,
		getDetailQryHdl: getDetailQryHdl,
		listQryHdl:      listQryHdl,
	}
}
