package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BillingCustomer struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey;column:id"`
	BarbershopID     uuid.UUID `gorm:"type:uuid;not null;column:barbershop_id"`
	StripeCustomerID string    `gorm:"type:text;not null;uniqueIndex;column:stripe_customer_id"`
	CreatedAt        time.Time `gorm:"type:timestamptz;not null;column:created_at"`
	UpdatedAt        time.Time `gorm:"type:timestamptz;not null;column:updated_at"`
}

func (BillingCustomer) TableName() string { return "billing_customers" }

func (bc *BillingCustomer) BeforeCreate(tx *gorm.DB) error {
	if bc.ID == uuid.Nil {
		bc.ID = uuid.New()
	}
	return nil
}
