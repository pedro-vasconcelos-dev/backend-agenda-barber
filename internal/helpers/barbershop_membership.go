package helpers

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// HasBarbershopAccess retorna true se o usuário pertence à barbearia com qualquer role ativo.
func HasBarbershopAccess(db *gorm.DB, userID string, barbershopID uuid.UUID) (bool, error) {
	return hasBarbershopRole(db, userID, barbershopID, "")
}

// IsBarbershopOwner retorna true se o usuário é owner ativo da barbearia.
func IsBarbershopOwner(db *gorm.DB, userID string, barbershopID uuid.UUID) (bool, error) {
	return hasBarbershopRole(db, userID, barbershopID, "owner")
}

func hasBarbershopRole(db *gorm.DB, userID string, barbershopID uuid.UUID, role string) (bool, error) {
	q := db.Table("barbershop_users").
		Where("user_id = ? AND barbershop_id = ? AND is_active = true", userID, barbershopID)
	if role != "" {
		q = q.Where("role = ?", role)
	}

	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
