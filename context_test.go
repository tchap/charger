package charger_test

import (
	"testing"

	"github.com/tchap/charger"
)

func TestContext_RenderValues_ErrCyrcleDetected(t *testing.T) {
	ch := charger.New()

	ch.AddKey(charger.String{
		Name:    "A",
		Default: `{{ get "B" }}`,
	})

	ch.AddKey(charger.String{
		Name:    "B",
		Default: `{{ get "A" }}`,
	})

	if _, err := ch.GatherValues(); err != ErrCyrcleDetected {
		t.Error("expected to detect a dependency cyrcle")
	}
}
