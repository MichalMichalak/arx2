package idgen

import (
	"github.com/google/uuid"
)

type IDGenerator interface {
	NewId() string
}

type UUIDGenerator struct{}

func (UUIDGenerator) NewId() string {
	return uuid.New().String()
}
