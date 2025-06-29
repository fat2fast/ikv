package userservice

import (
	"context"
	"errors"

	usermodel "fat2fast/ikv/modules/user/model"
	"fat2fast/ikv/shared"
	"fat2fast/ikv/shared/datatype"
)

type AuthenticateCommand struct {
	Dto usermodel.LoginForm
}

type IAuthenticateRepo interface {
	FindByEmail(ctx context.Context, email string) (*usermodel.User, error)
}
type ITokenIssuer interface {
	IssueToken(ctx context.Context, userID string) (string, error)
	ExpIn() int
}

type AuthenticateCommandHandler struct {
	repo        IAuthenticateRepo
	tokenIssuer ITokenIssuer
}
type AuthenticateResult struct {
	Token string `json:"token"`
	ExpIn int    `json:"expIn"`
}

func NewAuthenticateCommandHandler(repo IAuthenticateRepo, tokenIssuer ITokenIssuer) *AuthenticateCommandHandler {
	return &AuthenticateCommandHandler{repo: repo, tokenIssuer: tokenIssuer}
}
func (hdl *AuthenticateCommandHandler) Execute(ctx context.Context, cmd *AuthenticateCommand) (*AuthenticateResult, error) {

	user, err := hdl.repo.FindByEmail(ctx, cmd.Dto.Username)
	if err != nil {
		if errors.Is(err, datatype.ErrRecordNotFound) {
			return nil, datatype.ErrBadRequest.WithError(usermodel.ErrInvalidEmailAndPassword.Error())
		}

		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	if user.Status == usermodel.StatusDeleted || user.Status == usermodel.StatusBanned {
		return nil, datatype.ErrBadRequest.WithError(usermodel.ErrUserBannedOrDeleted.Error())
	}
	// Kiểm tra password bằng bcrypt.CompareHashAndPassword
	err = shared.VerifyPassword(cmd.Dto.Password, user.Salt, user.Password)
	if err != nil {
		return nil, datatype.ErrBadRequest.WithError(usermodel.ErrInvalidEmailAndPassword.Error())
	}
	token, err := hdl.tokenIssuer.IssueToken(ctx, user.ID.String())
	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}
	return &AuthenticateResult{Token: token, ExpIn: hdl.tokenIssuer.ExpIn()}, nil

}
