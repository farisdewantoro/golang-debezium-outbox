package order

import (
	"time"

	"github.com/go-openapi/strfmt"
	"gorm.io/gorm"
)

type Order struct {
	ID        strfmt.UUID4   `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	UserID    uint           `json:"user_id"`
	Total     float64        `json:"total"`
	Status    string         `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type CreateOrderParam struct {
	UserID uint    `json:"user_id"`
	Total  float64 `json:"total"`
}

func (param *CreateOrderParam) ToDomain() *Order {
	return &Order{
		UserID: param.UserID,
		Total:  param.Total,
		Status: "PENDING", // Initial status for new orders
	}
}
