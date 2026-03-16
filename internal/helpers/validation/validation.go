package validation

import (
	"errors"
	"math/rand"
	"regexp"
	"strings"
	"time"

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

func GenRandomPassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	if length <= 0 {
		return "", errors.New("length must be greater than 0")
	}

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	var password strings.Builder
	for i := 0; i < length; i++ {
		randomIndex := seededRand.Intn(len(charset))
		password.WriteByte(charset[randomIndex])
	}
	return password.String(), nil
}
