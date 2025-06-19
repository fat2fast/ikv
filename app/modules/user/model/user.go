package usermodel

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"column:id;"`
	Name      string    `json:"name" gorm:"column:name;"`
	Email     string    `json:"email" gorm:"column:email;"`
	Password  string    `json:"password" gorm:"column:password;"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;"`
	Status    string    `json:"status" gorm:"column:status;"`
}

const (
	StatusActive   = "active"
	StatusInactive = "inactive"
	StatusPending  = "pending"
	StatusDeleted  = "deleted"
)
