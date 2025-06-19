package userhttpgin

import (
	"context"

	usermodel "fat2fast/ikv/modules/user/model"
	usersevice "fat2fast/ikv/modules/user/service"

	"github.com/google/uuid"
)

type ICreateCommandHandler interface {
	Execute(ctx context.Context, cmd *usersevice.CreateCommand) (*uuid.UUID, error)
}
type IDetailQueryHandler interface {
	Execute(ctx context.Context, query *usersevice.GetDetailsQuery) (*usermodel.User, error)
}
type UserHTTPController struct {
	createCmdHdl    ICreateCommandHandler
	getDetailQryHdl IDetailQueryHandler
	// updateCmdHdl    IUpdateByIdCommandHandler
	// deleteCmdHdl    IDeleteByIdCommandHandler
	// listQryHdl      IListQueryHandler
}

func NewUserHTTPController(
	createCmdHdl ICreateCommandHandler,
	getDetailQryHdl IDetailQueryHandler,
	// updateCmdHdl IUpdateByIdCommandHandler,
	// deleteCmdHdl IDeleteByIdCommandHandler,
	// listQryHdl IListQueryHandler,
	// repoRPCCategory IRepoRPCCategory,
) *UserHTTPController {
	return &UserHTTPController{
		getDetailQryHdl: getDetailQryHdl,
		createCmdHdl:    createCmdHdl,
		// updateCmdHdl:    updateCmdHdl,
		// deleteCmdHdl:    deleteCmdHdl,
		// listQryHdl:      listQryHdl,
		// repoRPCCategory: repoRPCCategory,
	}
}
