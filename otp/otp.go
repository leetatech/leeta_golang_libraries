package otp

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/rs/zerolog/log"
)

// Generator defines the interface for generating OTP codes.
type Generator interface {
	Generate() (string, error)
}

type generator struct{}

// New creates and returns a new instance of the OTP generator.
func New() Generator {
	return &generator{}
}

// Generate creates a new 6-digit numeric OTP using cryptographic randomness.
func (o *generator) Generate() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		log.Error().Msg("unable to generate otp")
		return "", fmt.Errorf("unable to generate otp: %w", err)
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
