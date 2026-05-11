package servicerequest

import "time"

type ServiceRequest struct {
	ID                 int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CustomerName       string    `gorm:"column:customer_name;not null" json:"customer_name"`
	CustomerPhone      string    `gorm:"column:customer_phone;not null" json:"customer_phone"`
	DeviceMake         string    `gorm:"column:device_make" json:"device_make,omitempty"`
	DeviceModel        string    `gorm:"column:device_model" json:"device_model,omitempty"`
	ProblemDescription string    `gorm:"column:problem_description" json:"problem_description,omitempty"`
	Status             string    `gorm:"column:status;default:'new'" json:"status"`
	CreatedAt          time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (ServiceRequest) TableName() string { return "service_requests" }

type CreateServiceRequestInput struct {
	CustomerName       string `json:"customer_name" binding:"required"`
	CustomerPhone      string `json:"customer_phone" binding:"required"`
	DeviceMake         string `json:"device_make"`
	DeviceModel        string `json:"device_model"`
	ProblemDescription string `json:"problem_description"`
}

type UpdateStatusInput struct {
	Status string `json:"status" binding:"required"`
}
