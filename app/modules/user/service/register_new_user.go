package userservice

import (
	"context"
	"time"

	usermodel "fat2fast/ikv/modules/user/model"
	"fat2fast/ikv/shared"
	"fat2fast/ikv/shared/datatype"

	"github.com/google/uuid"
)

type CreateCommand struct {
	Dto usermodel.RegisterForm
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
func (uc *CreateCommandHandler) Execute(ctx context.Context, cmd *CreateCommand) (*usermodel.RegisterResponse, error) {

	salt, _ := shared.RandomStr(16)
	hashPassword, err := shared.HashPassword(cmd.Dto.Password, salt)

	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	newId, _ := uuid.NewV7()
	now := time.Now().In(time.FixedZone("Asia/Ho_Chi_Minh", 7*3600))

	user := &usermodel.User{
		ID:        newId,
		Email:     cmd.Dto.Email,
		Password:  string(hashPassword),
		FirstName: cmd.Dto.FirstName,
		LastName:  cmd.Dto.LastName,
		Phone:     cmd.Dto.Phone,
		Salt:      salt,
		Status:    usermodel.StatusActive,
		Type:      usermodel.TypeEmailPassword,
		Role:      usermodel.RoleUser,
		CreatedAt: &now,
		CreatedBy: "system",
		UpdatedAt: &now,
		UpdatedBy: "system",
	}
	err = uc.userRepo.Insert(ctx, user)
	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}
	userResponse := &usermodel.RegisterResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
	}
	return userResponse, nil
}
