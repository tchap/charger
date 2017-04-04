package charger

type Lookuper interface {
	Lookup(key string) (string, error)
}
