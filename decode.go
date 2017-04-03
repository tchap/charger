package charger

import (
	"encoding"
)

type Charger interface {
	Charge(*Charger, string) error
}

func Charge(ctx Context, v interface{}) error {
	if charger, ok := v.(Charger); ok {
		return charger.Charge(data)
	}

	if unmarshaler, ok := v.(encoding.TextUnmarshaler); ok {
		return unmarshaler.UnmarshalText(data)
	}

}
