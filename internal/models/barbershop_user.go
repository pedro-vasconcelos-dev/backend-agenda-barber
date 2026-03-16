package models

import (
	"time"

	"github.com/google/uuid"
)

// CREATE TABLE IF NOT EXISTS barbershop_users (
//   id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
//   barbershop_id uuid NOT NULL REFERENCES barbershops(id) ON DELETE CASCADE,
//   user_id       uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
//   role          text NOT NULL CHECK (role IN ('owner', 'reception', 'professional')),
//   is_active     boolean NOT NULL DEFAULT true,
//   created_at    timestamptz NOT NULL DEFAULT now(),
//   updated_at    timestamptz NOT NULL DEFAULT now(),
//   UNIQUE (barbershop_id, user_id)
// );

type BarbershopUser struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	BarbershopID uuid.UUID `gorm:"type:uuid;not null" json:"barbershop_id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Role         string    `gorm:"type:text;not null;check:role IN ('owner', 'reception', 'professional')" json:"role"`
	IsActive     bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt    time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt    time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

func (BarbershopUser) TableName() string {
	return "barbershop_users"
}
