package models

import (
	"time"

	"github.com/google/uuid"
)

// CREATE TABLE IF NOT EXISTS users (
//   id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
//   name          text NOT NULL,
//   email         text NOT NULL,
//   phone         text,
//   password_hash text NOT NULL,
//   is_admin      boolean NOT NULL DEFAULT false,
//   is_active     boolean NOT NULL DEFAULT true,
//   created_at    timestamptz NOT NULL DEFAULT now(),
//   updated_at    timestamptz NOT NULL DEFAULT now()
// );

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name         string    `gorm:"not null" json:"name"`
	Email        string    `gorm:"not null;unique" json:"email"`
	Phone        string    `gorm:"unique" json:"phone,omitempty"`
	PasswordHash string    `gorm:"not null" json:"-"`
	IsAdmin      bool      `gorm:"not null;default:false" json:"is_admin"`
	IsActive     bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt    time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt    time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
