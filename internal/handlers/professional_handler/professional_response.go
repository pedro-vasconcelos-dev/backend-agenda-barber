package professional_handler

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProfessionalResponse é o tipo de resposta compartilhado entre get, list e update.
type ProfessionalResponse struct {
	ID           string `json:"id"`
	BarbershopID string `json:"barbershop_id"`
	UserID       string `json:"user_id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Role         string `json:"role"`
	IsActive     bool   `json:"is_active"`
}

// professionalScan é o alvo do Scan do GORM para queries com join.
type professionalScan struct {
	ID           uuid.UUID `gorm:"column:id"`
	BarbershopID uuid.UUID `gorm:"column:barbershop_id"`
	UserID       uuid.UUID `gorm:"column:user_id"`
	Name         string    `gorm:"column:name"`
	Phone        string    `gorm:"column:phone"`
	IsActive     bool      `gorm:"column:is_active"`
	Email        string    `gorm:"column:email"`
	Role         string    `gorm:"column:role"`
}

func (r professionalScan) toResponse() ProfessionalResponse {
	return ProfessionalResponse{
		ID:           r.ID.String(),
		BarbershopID: r.BarbershopID.String(),
		UserID:       r.UserID.String(),
		Name:         r.Name,
		Email:        r.Email,
		Phone:        r.Phone,
		Role:         r.Role,
		IsActive:     r.IsActive,
	}
}

// professionalQuery retorna a base de query com joins em users e barbershop_users.
func professionalQuery(db *gorm.DB) *gorm.DB {
	return db.Table("professionals p").
		Select("p.id, p.barbershop_id, p.user_id, p.name, p.phone, p.is_active, u.email, bu.role").
		Joins("JOIN users u ON u.id = p.user_id").
		Joins("JOIN barbershop_users bu ON bu.user_id = p.user_id AND bu.barbershop_id = p.barbershop_id")
}
