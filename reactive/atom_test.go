package reactive

import (
	"testing"

	"github.com/influx6/flux"
)

func TestStrictAtom(t *testing.T) {
	models := StrictAtom("model", false)

	m, ok := models.Mutate("admin")
	if !ok && m.Value() == "admin" {
		flux.FatalFailed(t, "Unable to mutate value")
	}

	flux.LogPassed(t, "StrictAtom: Value mutate to 'admin' successfully")

	m, ok = models.Mutate(1)
	if ok && m.Value() == 1 {
		flux.FatalFailed(t, "Mutate value to incorrect type:", m.Value())
	}

	flux.LogPassed(t, "StrictAtom: Value not mutated to 1 with atom type set to string!")

}

func TestUnstrictAtom(t *testing.T) {
	models := UnstrictAtom("model", false)

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
