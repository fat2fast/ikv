package shared

import (
	"encoding/hex"
	"fmt"
	"math/rand"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func RandomStr(length int) (string, error) {
	var b = make([]byte, length)

	_, err := rand.Read(b)

	if err != nil {
		return "", errors.WithStack(err)
	}

	return hex.EncodeToString(b), nil
}

func HashPassword(password string, salt string) (string, error) {
	saltPass := fmt.Sprintf("%s.%s", salt, password)
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(saltPass), 8)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return string(hashPassword), nil
}

// VerifyPassword kiểm tra password có khớp với hash đã lưu không
func VerifyPassword(password string, salt string, hashedPassword string) error {
	saltPass := fmt.Sprintf("%s.%s", salt, password)
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(saltPass))
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
