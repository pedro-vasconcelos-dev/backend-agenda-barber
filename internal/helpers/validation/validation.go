package validation

import (
	"errors"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// NormalizeEmail normaliza o email para lowercase e remove espaços
func NormalizeEmail(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// NormalizePhone remove todos os caracteres não numéricos do número de telefone
// ex: "+55 (17) 988299693" -> "5517988299693"
func NormalizePhone(phone string) string {
	if phone == "" {
		return ""
	}
	// Remove todos os caracteres que não são dígitos
	re := regexp.MustCompile(`[^0-9]`)
	return re.ReplaceAllString(phone, "")
}

// HashPassword gera hash da senha usando bcrypt
func HashPassword(plain string) (string, error) {
	if len(plain) < 6 {
		return "", errors.New("password must be at least 6 characters")
	}
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// CheckPassword verifica se a senha corresponde ao hash
func CheckPassword(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
