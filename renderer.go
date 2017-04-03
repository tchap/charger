package charger

type Getter interface {
	Get(string) string
}

type GetterFunc func(string) string

func (fnc GetterFunc) Get(key string) string {
	return fnc(key)
}

type Renderer interface {
	Dependencies(string) ([]string, error)
	Render(string, Getter) (string, error)
}
