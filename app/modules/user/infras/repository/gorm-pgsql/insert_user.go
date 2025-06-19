package userrepository

import (
	"context"

	usermodel "fat2fast/ikv/modules/user/model"

	"github.com/pkg/errors"
)

func (repo *UserRepository) Insert(ctx context.Context, data *usermodel.User) error {
	db := repo.dbCtx.GetMainConnection()

	if err := db.Create(data).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
