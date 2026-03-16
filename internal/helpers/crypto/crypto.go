package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/bcrypt"
)

var encryptionKey = []byte("z4T9mY2pL8cNwQ1rX6jFvU0sB3eH5uK!") // agora 32 bytes

// GenerateRandomPassword retorna uma senha aleatória do tamanho especificado.
func GenerateRandomPassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b), nil
}

// HashPassword gera um hash bcrypt para autenticação.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword compara hash bcrypt com a senha em texto.
func CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// Encrypt cifra o texto (AES-CFB) e retorna base64(IV|ciphertext).
func Encrypt(plainText string) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(plainText))
	stream.XORKeyStream(ciphertext, []byte(plainText))
	combined := append(iv, ciphertext...)
	return base64.StdEncoding.EncodeToString(combined), nil
}

// Decrypt decifra um texto cifrado por Encrypt.
func Decrypt(cipherText string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	if len(data) < aes.BlockSize {
		return "", errors.New("ciphertext muito curto")
	}
	iv, ct := data[:aes.BlockSize], data[aes.BlockSize:]
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ct, ct)
	return string(ct), nil
}
