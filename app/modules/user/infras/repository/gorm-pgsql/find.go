package userrepository

import (
	"context"

	usermodel "fat2fast/ikv/modules/user/model"
	"fat2fast/ikv/shared/datatype"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (*usermodel.User, error) {
	return repo.FindByCondition(ctx, map[string]interface{}{"email": email})
}

func (repo *UserRepository) FindById(ctx context.Context, id uuid.UUID) (*usermodel.User, error) {
	return repo.FindByCondition(ctx, map[string]interface{}{"id": id})
}

func (repo *UserRepository) FindByCondition(ctx context.Context, cond map[string]interface{}) (*usermodel.User, error) {
	var user usermodel.User

	db := repo.dbCtx.GetMainConnection()

	if err := db.Table(user.TableName()).Where(cond).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, datatype.ErrRecordNotFound
		}

		return nil, errors.WithStack(err)
	}

	return &user, nil
}
