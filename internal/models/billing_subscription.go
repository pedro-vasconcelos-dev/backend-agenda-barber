package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BillingSubscription struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;column:id"`

	BarbershopID uuid.UUID `gorm:"type:uuid;not null;column:barbershop_id"`

	StripeCustomerID     string `gorm:"type:text;not null;column:stripe_customer_id"`
	StripeSubscriptionID string `gorm:"type:text;not null;uniqueIndex;column:stripe_subscription_id"`

	StripePriceLookupKey string `gorm:"type:text;not null;column:stripe_price_lookup_key"`

	Status            string     `gorm:"type:text;not null;column:status"`
	CurrentPeriodEnd  *time.Time `gorm:"column:current_period_end"`
	CancelAtPeriodEnd bool       `gorm:"column:cancel_at_period_end"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (BillingSubscription) TableName() string { return "billing_subscriptions" }

func (s *BillingSubscription) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
