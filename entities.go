package bookclub

import (
	"errors"
	"net/mail"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
)

var (
	ErrPasswordMismatch   error = errors.New("password mismatch")
	ErrInvalidPassword          = errors.New("password must be between 8 and 30 characters and contain at least one number, one uppercase letter, one lowercase letter, and one special character")
	ErrInvalidEmail       error = errors.New("invalid email")
	ErrInternalError      error = errors.New("internal error")
	ErrAlreadyExists            = errors.New("resource already exists")
	ErrInvalidCredentials       = errors.New("invalid credentials")
	ErrNotFound                 = errors.New("resource not found")
)

type User struct {
	ID        uuid.UUID
	Email     Email
	password  Password
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (o *User) SetPassword(password Password) {
	o.password = password
}

func (o *User) ComparePassword(password string) error {
	return o.password.Compare(password, o.password)
}

type Password []byte

func NewPassword(password string) (Password, error) {
	if !isValidPassword(password) {
		return nil, ErrInvalidPassword
	}

	ans, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, ErrInternalError
	}

	return ans, nil
}

func (o Password) Compare(password string, hashedPassword Password) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return ErrPasswordMismatch
	}

	return nil
}

type Email string

func NewEmail(value string) (Email, error) {
	if _, err := mail.ParseAddress(value); err != nil {
		return "", ErrInvalidEmail
	}

	return Email(value), nil
}

func isValidPassword(password string) bool {
	if len(password) < 8 || len(password) > 30 {
		return false
	}

	var (
		hasLower   bool
		hasUpper   bool
		hasDigit   bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasLower && hasUpper && hasDigit && hasSpecial
}
