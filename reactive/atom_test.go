package reactive

import "testing"

func TestStrictAtom(t *testing.T) {
	models := StrictAtom("model", false)

	m, ok := models.Mutate("admin")
	if !ok && m.Value() == "admin" {
		t.Fatal("Unable to mutate value")
	}

	m, ok = models.Mutate(1)
	if ok && m.Value() == 1 {
		t.Fatal("mutate value to incorrect type:", m.Value())
	}

}

func TestUnstrictAtom(t *testing.T) {
	models := UnstrictAtom("model", false)

	m, ok := models.Mutate("admin")

	if !ok && m.Value() == "admin" {
		t.Fatal("Unable to mutate value")
	}

	m, ok = models.Mutate(1)

	if !ok && m.Value() != 1 {
		t.Fatal("Unable to mutate value to int type", m.Value())
	}
}
