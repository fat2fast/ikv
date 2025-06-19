package userrepository

import (
	"context"

	usermodel "fat2fast/ikv/modules/user/model"
	"fat2fast/ikv/shared/datatype"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (repo *UserRepository) FindById(ctx context.Context, id uuid.UUID) (*usermodel.User, error) {
	var user usermodel.User

	db := repo.dbCtx.GetMainConnection()

	if err := db.Where("id = ?", id.String()).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, datatype.ErrRecordNotFound
		}
		return nil, errors.WithStack(err)
	}

	return &user, nil
}
