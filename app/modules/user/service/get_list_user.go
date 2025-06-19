package userservice

import (
	"context"
	usermodel "fat2fast/ikv/modules/user/model"
)

type ListQuery struct {
	Page  int
	Limit int
}

type IListRepo interface {
	List(ctx context.Context, query *ListQuery) ([]*usermodel.User, error)
}
