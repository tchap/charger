package charger

import (
	"errors"
)

var ErrKeyTaken = errors.New("key already taken")

type ErrRequired struct {
	Name string
}

func (err *ErrRequired) Error() string {
	return "key required but not set: " + err.Name
}

type Lookuper interface {
	Lookup(key string) (string, error)
}

type Key interface {
	GetName() string
	GetDefault() (string, bool)
	GetRequired() bool

	GetLoaders() []Loader
}

type keyRecord struct {
	Key   Key
	Value string
}

type Charger struct {
	keys map[string]*keyRecord
}

func (ch *Charger) AddKey(key Key) error {
	name := key.GetName()

	if _, ok := ch.records[name]; ok {
		return ErrKeyTaken
	}

	ch.records[name] = &keyRecord{key, ""}
	return nil
}

func (ch *Charger) MustAddKey(key Key) {
	if err := ch.AddKey(key); err != nil {
		panic(err)
	}
}

func (ch *Charger) GatherValues() (*Context, error) {
	ctx := &Context{
		keys: make(map[string]*keyRecord, len(ch.keys)),
	}

KeyLoop:
	for name, record := range ch.keys {
		var value string

		for _, lookuper := range record.Key.Lookupers() {
			value, err = lookuper.Lookup(name)
			switch err {
			case nil:
				ctx.keys[name] = &keyRecord{record.Key, value}
				continue KeyLoop
			case ErrNotFound:
			default:
				return nil, err
			}
		}

		for _, lookuper := range ch.lookupers {
			value, err = lookuper.Lookup(name)
			switch err {
			case nil:
				ctx.keys[name] = &keyRecord{record.Key, value}
				continue KeyLoop
			case ErrNotFound:
			default:
				return nil, err
			}
		}

		value, ok := record.Key.GetDefault()
		if ok {
			ctx.keys[name] = &keyRecord{record.Key, value}
			continue KeyLoop
		}

		if record.Key.GetRequired() {
			return nil, &ErrRequired{name}
		}
	}

	if err := ctx.render(); err != nil {
		return nil, err
	}

	return ctx, nil
}

type Renderer interface {
	Dependencies() ([]string, error)
	Render(*Context, string) error
}

func (ctx *Context) render() error {
	deps := make(map[string]*[]string, len(ctx.keys))

	for name, record := range ctx.keys {
		depList, err := record.Key.Renderer().Dependencies()
		if err != nil {
			return err
		}

		deps[name] = &depList
	}

	fulfilled := make(map[string]struct{}, len(ctx.keys))

	for {
		if len(fulfilled) == len(ctx.keys) {
			return nil
		}

		var removed bool

	DepLoop:
		for name, depList := range deps {
			if _, ok := fulfilled[name]; ok {
				continue
			}

			for _, dep := range depList {
				if _, ok := fulfilled[dep]; !ok {
					continue DepLoop
				}
			}

			fulfilled[name] = struct{}{}
			removed = true
		}

		if !removed {
			return ErrCyrcleDetected
		}
	}
}
