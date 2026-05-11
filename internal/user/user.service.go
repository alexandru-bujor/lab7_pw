package user

import (
	"MegaMobileBack/internal/utils"
	"MegaMobileBack/pkg/db"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func Login(email, password string) (User, error) {
	var user User

	result := db.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return user, errors.New("user not found")
	}

	if !utils.CheckPasswordSHA256(password, user.Password) {
		return user, errors.New("invalid password")
	}

	return user, nil
}

type Service struct {
	DB *gorm.DB
}

func NewService() *Service {
	return &Service{DB: db.DB}
}

func (s *Service) List() ([]User, error) {
	var users []User
	if err := s.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Service) GetByID(id int) (*User, error) {
	var user User
	if err := s.DB.First(&user, id).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (s *Service) GetByEmail(email string) (*User, error) {
	var user User
	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (s *Service) Create(input CreateUserInput) (*User, error) {
	existing, _ := s.GetByEmail(input.Email)
	if existing != nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword := utils.HashPasswordSHA256(input.Password)
	permissionsJSON, err := json.Marshal(input.Permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal permissions: %v", err)
	}

	user := User{
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		Email:       input.Email,
		Password:    hashedPassword,
		Role:        input.Role,
		Permissions: permissionsJSON,
	}

	if err := s.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) Update(id int, input UpdateUserInput) (*User, error) {
	var user User
	if err := s.DB.First(&user, id).Error; err != nil {
		return nil, errors.New("user not found")
	}

	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}
	if input.LastName != "" {
		user.LastName = input.LastName
	}
	if input.Email != "" && input.Email != user.Email {
		existing, _ := s.GetByEmail(input.Email)
		if existing != nil && existing.ID != id {
			return nil, errors.New("user with this email already exists")
		}
		user.Email = input.Email
	}
	if input.Password != "" {
		user.Password = utils.HashPasswordSHA256(input.Password)
	}
	if input.Role != "" {
		user.Role = input.Role
	}
	if input.Permissions != nil {
		permissionsJSON, err := json.Marshal(input.Permissions)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal permissions: %v", err)
		}
		user.Permissions = permissionsJSON
	}

	if err := s.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) Delete(id int) error {
	return s.DB.Delete(&User{}, id).Error
}
