package sharedinfras

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IDbContext interface {
	GetMainConnection() *gorm.DB
}
type IMiddlewareProvider interface {
	Auth() gin.HandlerFunc
	CheckRoles(roles ...string) gin.HandlerFunc
}

type IAppContext interface {
	DbContext() IDbContext
	MiddlewareProvider() IMiddlewareProvider
	// GetConfig() *datatype.Config
	// Uploader() IUploader
	// MsgBroker() IMsgBroker
}

type appContext struct {
	dbContext   IDbContext
	mldProvider IMiddlewareProvider
	// config      *datatype.Config
	// uploader    IUploader
	// msgBroker   IMsgBroker
}

func NewAppContext(db *gorm.DB) IAppContext {
	dbCtx := NewDbContext(db)

	// introspectRpcClient := sharedrpc.NewIntrospectRpcClient(config.UserServiceURL)

	// provider := middleware.NewMiddlewareProvider(introspectRpcClient)

	return &appContext{
		dbContext: dbCtx,
		// mldProvider: provider,
	}
}

func (c *appContext) MiddlewareProvider() IMiddlewareProvider {
	return c.mldProvider
}

func (c *appContext) DbContext() IDbContext {
	return c.dbContext
}
