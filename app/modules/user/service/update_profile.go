package userservice

import (
	"context"
	"time"

	usermodel "fat2fast/ikv/modules/user/model"
	"fat2fast/ikv/shared/datatype"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// UpdateProfileCommand đại diện cho command cập nhật profile user
type UpdateProfileCommand struct {
	UserID uuid.UUID
	Dto    usermodel.UpdateProfileRequest
}

// IUpdateProfileRepo interface cho repository operations cần thiết
type IUpdateProfileRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (*usermodel.User, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error
}

// UpdateProfileCommandHandler xử lý command cập nhật profile
type UpdateProfileCommandHandler struct {
	repo IUpdateProfileRepo
}

// NewUpdateProfileCommandHandler khởi tạo handler mới
func NewUpdateProfileCommandHandler(repo IUpdateProfileRepo) *UpdateProfileCommandHandler {
	return &UpdateProfileCommandHandler{repo: repo}
}

// Execute thực thi command cập nhật profile
func (hdl *UpdateProfileCommandHandler) Execute(ctx context.Context, cmd *UpdateProfileCommand) error {
	// Validate input
	if cmd.UserID == uuid.Nil {
		return datatype.ErrBadRequest.WithError("User ID is required")
	}

	// Kiểm tra user có tồn tại không
	user, err := hdl.repo.FindById(ctx, cmd.UserID)
	if err != nil {
		if errors.Is(err, datatype.ErrRecordNotFound) {
			return datatype.ErrNotFound.WithError("User not found")
		}
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// Kiểm tra trạng thái user
	if user.Status == usermodel.StatusDeleted {
		return datatype.ErrDeleted.WithError("Cannot update deleted user")
	}

	if user.Status == usermodel.StatusBanned {
		return datatype.ErrForbidden.WithError("Cannot update banned user")
	}

	// Chuẩn bị dữ liệu cập nhật
	updates := make(map[string]interface{})

	if cmd.Dto.FirstName != "" {
		updates["first_name"] = cmd.Dto.FirstName
	}

	if cmd.Dto.LastName != "" {
		updates["last_name"] = cmd.Dto.LastName
	}

	if cmd.Dto.Phone != "" {
		updates["phone"] = cmd.Dto.Phone
	}

	// Thêm thông tin audit
	updates["updated_at"] = time.Now()
	updates["updated_by"] = cmd.UserID.String()

	// Kiểm tra có thay đổi gì không
	if len(updates) <= 2 { // Chỉ có updated_at và updated_by
		return datatype.ErrBadRequest.WithError("No fields to update")
	}

	// Thực hiện cập nhật
	if err := hdl.repo.UpdateProfile(ctx, cmd.UserID, updates); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug("Failed to update user profile")
	}

	return nil
}
