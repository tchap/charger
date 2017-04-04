package charger

import (
	"fmt"
)

type Key interface {
	Name() string
	Default() (string, bool)
	Required() bool

	Lookupers() []Lookuper
	Renderer() Renderer
}

func newDomainKey(key Key, domain string) Key {
	return &domainKey{key, domain}
}

type domainKey struct {
	Key
	domain string
}

func (key *domainKey) Name() string {
	return fmt.Sprintf("%v.%v", key.domain, key.Key.Name())
}
