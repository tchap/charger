package charger_test

import (
	"reflect"
	"testing"

	"github.com/tchap/charger"
)

func TestTemplateRenderer_Dependencies(t *testing.T) {
	renderer := charger.NewTemplateRenderer()

	cases := []struct {
		tmpl string
		deps []string
	}{
		{`no deps at all`, nil},
		{`no deps at all`, []string{}},
		{`{{ get "A" }}`, []string{"A"}},
		{`{{ get "A" }}-{{ get "B" }}`, []string{"A", "B"}},
	}

	for i, c := range cases {
		// Get dependencies.
		deps, err := renderer.Dependencies(c.tmpl)
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
			t.Errorf("case %v: dependency list mismatch; expected %#v, got %#v", i, c.deps, deps)
		}
	}
}
