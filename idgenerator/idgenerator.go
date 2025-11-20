package idgenerator

import (
	"github.com/google/uuid"
)

// Generator defines the interface for generating unique IDs.
type Generator interface {
	Generate() string
}

// idGenerator is a concrete implementation of the Generator interface.
type idGenerator struct{}

// New creates and returns a new instance of the ID generator.
func New() Generator {
	return &idGenerator{}
}

// Generate creates a new UUID (Universally Unique Identifier) as a string.
func (g *idGenerator) Generate() string {
	return uuid.NewString()
}
