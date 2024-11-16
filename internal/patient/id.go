package patient

import "github.com/google/uuid"

// NewID generates a new ID
func NewID() ID {
	return ID{uuid.New()}
}

// ID represents a patient aggregate ID
type ID struct {
	uuid.UUID
}

// ParseID parses ID from string
func ParseID(id string) ID {
	return ID{uuid.MustParse(id)}
}
