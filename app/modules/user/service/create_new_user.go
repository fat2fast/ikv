package userservice

import (
	"context"

	usermodel "fat2fast/ikv/modules/user/model"

	"github.com/google/uuid"
)

type CreateCommand struct {
	Dto usermodel.User
}
type ICreateRepo interface {
	Insert(ctx context.Context, data *usermodel.User) error
}
type CreateCommandHandler struct {
	userRepo ICreateRepo
}

func NewCreateCommandHandler(userRepo ICreateRepo) *CreateCommandHandler {
	return &CreateCommandHandler{userRepo: userRepo}
}
func (uc *CreateCommandHandler) Execute(ctx context.Context, cmd *CreateCommand) (*uuid.UUID, error) {

	return nil, nil
}
