package xulid

import (
	"crypto/rand"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

// ULIDGenerator is a singleton struct for generating ULIDs.
type ULIDGenerator struct {
	entropy *ulid.MonotonicEntropy
}

var instance *ULIDGenerator
var once sync.Once

// Instance returns the singleton instance of ULIDGenerator.
func Instance() *ULIDGenerator {
	once.Do(func() {
		instance = &ULIDGenerator{
			entropy: ulid.Monotonic(rand.Reader, 0),
		}
	})
	return instance
}

// NewULID generates a new ULID.
func (gen *ULIDGenerator) NewULID() ULID {
	return ConvertULIDToXULID(ulid.MustNew(ulid.Timestamp(time.Now()), gen.entropy))
}
