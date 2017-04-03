package charger

import (
	"bytes"
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

func (renderer *TemplateRenderer) Dependencies(value string) ([]string, error) {
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

	t, err := template.New("").Funcs(funcMap).Parse(value)
	if err != nil {
		return nil, err
	}

	if err := t.Execute(ioutil.Discard, nil); err != nil {
		return nil, err
	}

	return deps, nil
}

func (renderer *TemplateRenderer) Render(value string, getter Getter) (string, error) {
	renderer.FuncMap["get"] = func(key string) string {
		return getter.Get(key)
	}

	t, err := template.New("").Funcs(renderer.FuncMap).Parse(value)
	if err != nil {
		return "", err
	}

	var out bytes.Buffer
	if err := t.Execute(&out, nil); err != nil {
		return "", err
	}

	return out.String(), nil
}
