package user

import (
	"encoding/json"
	"gorm.io/datatypes"
	"time"
)

type User struct {
	ID          int            `gorm:"column=id;primaryKey;autoIncrement" json:"id"`
	FirstName   string         `gorm:"column=first_name" json:"first_name"`
	LastName    string         `gorm:"column=last_name" json:"last_name"`
	Email       string         `gorm:"column=email;unique" json:"email"`
	Password    string         `gorm:"column=password" json:"-"`
	Role        string         `gorm:"column=role" json:"role"`
	Permissions datatypes.JSON `gorm:"column=permissions;type:json" json:"-"`
	CreatedAt   time.Time      `gorm:"column=created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column=updated_at" json:"updated_at"`
}

type UserResponse struct {
	ID          int       `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (u *User) ToResponse() UserResponse {
	var permissions []string
	if len(u.Permissions) > 0 {
		json.Unmarshal(u.Permissions, &permissions)
	}
	return UserResponse{
		ID:          u.ID,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Email:       u.Email,
		Role:        u.Role,
		Permissions: permissions,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

func (User) TableName() string {
	return "users"
}

type CreateUserInput struct {
	FirstName   string   `json:"first_name" binding:"required"`
	LastName    string   `json:"last_name" binding:"required"`
	Email       string   `json:"email" binding:"required,email"`
	Password    string   `json:"password" binding:"required,min=6"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

type UpdateUserInput struct {
	FirstName   string   `json:"first_name"`
	LastName    string   `json:"last_name"`
	Email       string   `json:"email" binding:"omitempty,email"`
	Password    string   `json:"password" binding:"omitempty,min=6"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}
