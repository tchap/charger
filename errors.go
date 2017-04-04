package charger

import (
	"errors"
)

var (
	ErrKeyTaken       = errors.New("key already taken")
	ErrCyrcleDetected = errors.New("cyrcle detected")
	ErrNotFound       = errors.New("not found")
)

type ErrRequired struct {
	Name string
}

func (err *ErrRequired) Error() string {
	return "key required but not set: " + err.Name
}
