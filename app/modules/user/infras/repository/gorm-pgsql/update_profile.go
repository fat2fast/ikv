package userrepository

import (
	"context"

	usermodel "fat2fast/ikv/modules/user/model"
	"fat2fast/ikv/shared/datatype"

	"github.com/google/uuid"
)

// UpdateProfile cập nhật thông tin profile của user
func (repo *UserRepository) UpdateProfile(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error {
	db := repo.dbCtx.GetMainConnection()

	// Thực hiện cập nhật với điều kiện user phải tồn tại và chưa bị xóa
	result := db.WithContext(ctx).
		Model(&usermodel.User{}).
		Where("id = ? AND status != ?", userID, usermodel.StatusDeleted).
		Updates(updates)

	if result.Error != nil {
		return datatype.ErrInternalServerError.WithWrap(result.Error).WithDebug("Failed to update user profile in database")
	}

	// Kiểm tra có rows nào được cập nhật không
	if result.RowsAffected == 0 {
		return datatype.ErrNotFound.WithError("User not found or already deleted")
	}

	return nil
}
