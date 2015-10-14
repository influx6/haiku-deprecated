package reactive

import (
	"testing"

	"github.com/influx6/flux"
)

func TestAtom(t *testing.T) {
	models := Atom("model", false)

	m, ok := models.Mutate("admin")

	if !ok && m.Value() == "admin" {
		flux.FatalFailed(t, "Unable to mutate value")
	}

	flux.LogPassed(t, "UnStrictAtom: Value mutate to 'admin' successfully")

	m, ok = models.Mutate(1)

	if !ok && m.Value() != 1 {
		flux.FatalFailed(t, "Mutate value to incorrect type:", m.Value())
	}

	flux.LogPassed(t, "UnstrictAtom: Value mutate to 1 successfully")
}
