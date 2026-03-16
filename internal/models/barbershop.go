package models

import (
	"time"

	"github.com/google/uuid"
)

// CREATE TABLE IF NOT EXISTS barbershops (
//   id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
//   name        text NOT NULL,
//   slug        text NOT NULL,
//   timezone    text NOT NULL DEFAULT 'America/Sao_Paulo',
//   address     text,
//   is_active   boolean NOT NULL DEFAULT true,
//   created_at  timestamptz NOT NULL DEFAULT now(),
//   updated_at  timestamptz NOT NULL DEFAULT now()
// );

type Barbershop struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Slug      string    `gorm:"not null;unique" json:"slug"`
	Timezone  string    `gorm:"not null;default:'America/Sao_Paulo'" json:"timezone"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	IsActive  bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

func (Barbershop) TableName() string {
	return "barbershops"
}
