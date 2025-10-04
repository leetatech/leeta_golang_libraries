// Package encrypto provides utilities for password encryption, validation, and email verification.
package encrypto

import (
	"errors"
	"strings"
	"unicode"

	"github.com/badoux/checkmail"
	"github.com/leetatech/leeta_golang_libraries/errs"
	"golang.org/x/crypto/bcrypt"
)

type encryptorHandler struct{}

var _ Manager = &encryptorHandler{}

type Manager interface {
	ComparePasscode(passcode, hashedPasscode string) error
	Generate(passcode string) ([]byte, error)
	ValidatePasswordStrength(s string) error
	ValidateEmailFormat(email string) error
	ValidateDomain(email, leetaDomain string) error
}

// New creates and returns a new instance of the password encryptor manager.
func New() Manager {
	return &encryptorHandler{}
}

// Generate hashes the provided passcode using bcrypt and returns the hashed password.
func (e *encryptorHandler) Generate(passcode string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(passcode), bcrypt.DefaultCost)
}

// ComparePasscode compares a plain passcode with a hashed passcode using bcrypt.
// Returns nil if they match, or an error if they do not.
func (e *encryptorHandler) ComparePasscode(passcode, hashedPasscode string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPasscode), []byte(passcode))
}

// ValidatePasswordStrength checks if the password meets minimum strength requirements:
// at least 6 characters, contains uppercase, lowercase, digit, and special character.
func (e *encryptorHandler) ValidatePasswordStrength(password string) error {
	const minLen = 6
	var hasUpper, hasLower, hasNumber, hasSpecial bool

	if len(password) < minLen {
		return errs.Body(errs.PasswordValidationError, errors.New("password must be at least six characters long"))
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		default:
			// If the character doesn't match any of the above, it's invalid
			return errs.Body(errs.PasswordValidationError, errors.New("password contains invalid characters"))
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return errs.Body(errs.PasswordValidationError, errors.New("password must contain at least one uppercase letter, one lowercase letter, one digit, and one special character"))
	}

	return nil
}

// ValidateEmailFormat checks if the provided email has a valid format.
func (e *encryptorHandler) ValidateEmailFormat(email string) error {
	if err := checkmail.ValidateFormat(email); err != nil {
		return err
	}
	return nil
}

// ValidateDomain checks if the email's domain matches the provided domain string
// and validates the email host.
func (e *encryptorHandler) ValidateDomain(email, domainString string) error {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return errs.Body(errs.EmailFormatError, nil)
	}
	domain := parts[1]

	if err := checkmail.ValidateHost(email); err != nil {
		return errs.Body(errs.ValidEmailHostError, err)
	}

	if strings.EqualFold(domain, domainString) {
		return errs.Body(errs.ValidLeetaDomainError, nil)
	}

	return nil
}
