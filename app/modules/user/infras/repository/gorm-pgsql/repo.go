package userrepository

import (
	sharedinfras "fat2fast/ikv/shared/infras"
)

type UserRepository struct {
	dbCtx sharedinfras.IDbContext
}

func NewUserRepository(dbCtx sharedinfras.IDbContext) *UserRepository {
	return &UserRepository{dbCtx: dbCtx}
}
