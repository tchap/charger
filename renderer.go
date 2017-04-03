package charger

type Renderer interface {
	Dependencies(string) ([]string, error)
	Render(*Context, string) (string, error)
}
