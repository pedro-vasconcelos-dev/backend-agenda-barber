package models

import (
	"time"

	"github.com/google/uuid"
)

type Professional struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	BarbershopID uuid.UUID  `gorm:"type:uuid;not null" json:"barbershop_id"`
	UserID       *uuid.UUID `gorm:"type:uuid" json:"user_id"`
	Name         string     `gorm:"type:text;not null" json:"name"`
	Phone        string     `gorm:"type:text;not null" json:"phone"`
	IsActive     bool       `gorm:"not null;default:true" json:"is_active"`
	CreatedAt    time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"not null;default:now()" json:"updated_at"`
}

func (Professional) TableName() string {
	return "professionals"
}
