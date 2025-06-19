package userservice

import (
	"context"

	usermodel "fat2fast/ikv/modules/user/model"
	"fat2fast/ikv/shared/datatype"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type GetDetailsQuery struct {
	Id uuid.UUID
}
type IGetDetailsRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (*usermodel.User, error)
}

type GetDetailsQueryHandler struct {
	restRepo IGetDetailsRepo
}

func NewGetDetailsQueryHandler(restRepo IGetDetailsRepo) *GetDetailsQueryHandler {
	return &GetDetailsQueryHandler{restRepo: restRepo}
}
func (hdl *GetDetailsQueryHandler) Execute(ctx context.Context, query *GetDetailsQuery) (*usermodel.User, error) {
	user, err := hdl.restRepo.FindById(ctx, query.Id)

	if err != nil {
		if errors.Is(err, datatype.ErrRecordNotFound) {
			return nil, datatype.ErrNotFound
		}

		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	if user.Status == usermodel.StatusDeleted {
		return nil, datatype.ErrDeleted.WithError(usermodel.ErrUserIsDeleted.Error())
	}

	return user, nil
}
