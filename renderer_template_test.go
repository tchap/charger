package charger_test

import (
	"reflect"
	"testing"

	"github.com/tchap/charger"
)

func TestTemplateRenderer_Dependencies(t *testing.T) {
	renderer := charger.NewTemplateRenderer()

	cases := []struct {
		value string
		deps  []string
	}{
		{`no deps at all`, nil},
		{`no deps at all`, []string{}},
		{`{{ get "A" }}`, []string{"A"}},
		{`{{ get "A" }}-{{ get "B" }}`, []string{"A", "B"}},
	}

	for i, c := range cases {
		// Get dependencies.
		deps, err := renderer.Dependencies(c.value)
		if err != nil {
			t.Error(err)
			continue
		}

		// Make sure []string{} and nil are treated as equal.
		if len(deps) == 0 && len(c.deps) == 0 {
			continue
		}

		// Compare the dep lists.
		if !reflect.DeepEqual(deps, c.deps) {
			t.Errorf("case %v: dependency list mismatch; expected %#v, got %#v", i+1, c.deps, deps)
		}
	}
}

func TestTemplateRenderer_Render(t *testing.T) {
	renderer := charger.NewTemplateRenderer()

	cases := []struct {
		in  string
		out string
	}{
		{`no deps at all`, `no deps at all`},
		{`{{ get "A" }}`, `A-value`},
		{`{{ get "A" }}-{{ get "B" }}`, `A-value-B-value`},
		{`{{ get "C" }}`, ``},
		{`{{ get "C" }}-suffix`, `-suffix`},
	}

	getter := charger.GetterFunc(func(key string) string {
		switch key {
		case "A":
			return "A-value"
		case "B":
			return "B-value"
		default:
			return ""
		}
	})

	for i, c := range cases {
		// Render.
		out, err := renderer.Render(c.in, getter)
		if err != nil {
			t.Error(err)
			continue
		}

		// Check.
		if out != c.out {
			t.Errorf("case %v: renderer output mismatch; expected %v, got %v", i+1, c.out, out)
		}
	}
}
