package userservice

import (
	"context"

	usermodel "fat2fast/ikv/modules/user/model"
	"fat2fast/ikv/shared/datatype"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// GetProfileQuery đại diện cho query lấy thông tin profile user
type GetProfileQuery struct {
	UserID uuid.UUID
}

// IGetProfileRepo interface cho repository operations cần thiết
type IGetProfileRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (*usermodel.User, error)
}

// GetProfileQueryHandler xử lý query lấy thông tin profile
type GetProfileQueryHandler struct {
	repo IGetProfileRepo
}

// NewGetProfileQueryHandler khởi tạo handler mới
func NewGetProfileQueryHandler(repo IGetProfileRepo) *GetProfileQueryHandler {
	return &GetProfileQueryHandler{repo: repo}
}

// Execute thực thi query lấy thông tin profile
func (hdl *GetProfileQueryHandler) Execute(ctx context.Context, query *GetProfileQuery) (*usermodel.ProfileResponse, error) {
	// Validate input
	if query.UserID == uuid.Nil {
		return nil, datatype.ErrBadRequest.WithError("User ID is required")
	}

	// Lấy thông tin user từ database
	user, err := hdl.repo.FindById(ctx, query.UserID)
	if err != nil {
		if errors.Is(err, datatype.ErrRecordNotFound) {
			return nil, datatype.ErrNotFound.WithError("User not found")
		}
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// Kiểm tra trạng thái user
	if user.Status == usermodel.StatusDeleted {
		return nil, datatype.ErrDeleted.WithError("User has been deleted")
	}

	if user.Status == usermodel.StatusBanned {
		return nil, datatype.ErrForbidden.WithError("User has been banned")
	}

	// Convert entity sang DTO response
	profileResponse := user.ToProfileResponse()

	return profileResponse, nil
}
