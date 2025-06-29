package usermodel

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type UserType string
type UserRole string
type UserStatus string

const (
	TypeEmailPassword UserType   = "email_password"
	TypeFacebook      UserType   = "facebook"
	TypeGmail         UserType   = "gmail"
	RoleUser          UserRole   = "user"
	RoleAdmin         UserRole   = "admin"
	StatusPending     UserStatus = "pending"
	StatusActive      UserStatus = "active"
	StatusInactive    UserStatus = "inactive"
	StatusBanned      UserStatus = "banned"
	StatusDeleted     UserStatus = "deleted"
)

type User struct {
	ID        uuid.UUID  `json:"id" gorm:"column:id;"`
	CreatedBy string     `json:"created_by" gorm:"column:created_by;"`
	CreatedAt *time.Time `json:"created_at" gorm:"column:created_at;"`
	UpdatedBy string     `json:"updated_by" gorm:"column:updated_by;"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;"`
	Status    UserStatus `json:"status" gorm:"column:status;"`
	Type      UserType   `json:"type" gorm:"column:type;"`
	Role      UserRole   `json:"role" gorm:"column:role;"`
	FirstName string     `json:"first_name" gorm:"column:first_name;"`
	LastName  string     `json:"last_name" gorm:"column:last_name;"`
	Phone     string     `json:"phone" gorm:"column:phone;"`
	Email     string     `json:"email" gorm:"column:email;"`
	Password  string     `json:"password" gorm:"column:password;"`
	Salt      string     `json:"salt" gorm:"column:salt;"`
}

func (User) TableName() string {
	return "user_users"
}

func (u *User) GetFullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

// ToProfileResponse chuyển đổi User entity sang ProfileResponse DTO
func (u *User) ToProfileResponse() *ProfileResponse {
	return &ProfileResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Phone:     u.Phone,
		FullName:  u.GetFullName(),
		Role:      u.Role,
		Status:    u.Status,
		Type:      u.Type,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
