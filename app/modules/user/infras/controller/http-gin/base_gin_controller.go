package userhttpgin

import (
	"context"

	usermodel "fat2fast/ikv/modules/user/model"
	usersevice "fat2fast/ikv/modules/user/service"
)

type ICreateCommandHandler interface {
	Execute(ctx context.Context, cmd *usersevice.CreateCommand) (*usermodel.RegisterResponse, error)
}
type IAuthenticateCommandHandler interface {
	Execute(ctx context.Context, cmd *usersevice.AuthenticateCommand) (*usersevice.AuthenticateResult, error)
}
type IGetProfileQueryHandler interface {
	Execute(ctx context.Context, query *usersevice.GetProfileQuery) (*usermodel.ProfileResponse, error)
}
type IUpdateProfileCommandHandler interface {
	Execute(ctx context.Context, cmd *usersevice.UpdateProfileCommand) error
}

type UserHTTPController struct {
	createCmdHdl        ICreateCommandHandler
	authenticateCmdHdl  IAuthenticateCommandHandler
	getProfileQryHdl    IGetProfileQueryHandler
	updateProfileCmdHdl IUpdateProfileCommandHandler
	// updateCmdHdl    IUpdateByIdCommandHandler
	// deleteCmdHdl    IDeleteByIdCommandHandler
	// listQryHdl      IListQueryHandler
}

func NewUserHTTPController(
	createCmdHdl ICreateCommandHandler,
	authenticateCmdHdl IAuthenticateCommandHandler,
	getProfileQryHdl IGetProfileQueryHandler,
	updateProfileCmdHdl IUpdateProfileCommandHandler,
	// updateCmdHdl IUpdateByIdCommandHandler,
	// deleteCmdHdl IDeleteByIdCommandHandler,
	// listQryHdl IListQueryHandler,
	// repoRPCCategory IRepoRPCCategory,
) *UserHTTPController {
	return &UserHTTPController{
		createCmdHdl:        createCmdHdl,
		authenticateCmdHdl:  authenticateCmdHdl,
		getProfileQryHdl:    getProfileQryHdl,
		updateProfileCmdHdl: updateProfileCmdHdl,
		// updateCmdHdl:    updateCmdHdl,
		// deleteCmdHdl:    deleteCmdHdl,
		// listQryHdl:      listQryHdl,
		// repoRPCCategory: repoRPCCategory,
	}
}
