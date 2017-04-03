package charger

import (
	"io/ioutil"
	"text/template"
)

type TemplateRenderer struct {
	FuncMap template.FuncMap
}

func NewTemplateRenderer() *TemplateRenderer {
	return &TemplateRenderer{
		FuncMap: make(template.FuncMap),
	}
}

func (renderer *TemplateRenderer) Dependencies(tmpl string) ([]string, error) {
	var deps []string

	funcMap := map[string]interface{}{
		"get": func(key string) string {
			for _, dep := range deps {
				if dep == key {
					return ""
				}
			}
			deps = append(deps, key)
			return ""
		},
	}

	t, err := template.New("").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return nil, err
	}

	if err := t.Execute(ioutil.Discard, nil); err != nil {
		return nil, err
	}

	return deps, nil
}

func (renderer *TemplateRenderer) Render(ctx *Context, tmpl string) (string, error) {
	panic("not implemented")
}
