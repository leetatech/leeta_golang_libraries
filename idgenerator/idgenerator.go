package idgenerator

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
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

// Generate creates a new ULID (Universally Unique Lexicographically Sortable Identifier) as a string.
func (generator *idGenerator) Generate() string {
	entropy := ulid.Monotonic(rand.Reader, 0)
	return ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
}
