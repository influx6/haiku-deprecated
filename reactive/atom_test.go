package reactive

import (
	"testing"
	"time"
)

func TestAtom(t *testing.T) {
	models := StrictAtom("model", true, nil)

	m, ok := models.Mutate("admin")
	if !ok && m.Value() == "admin" {
		t.Fatal("Unable to mutate value")
	}

	m, ok = models.Mutate(1)
	if ok && m.Value() != 1 {
		t.Fatal("Unable to mutate value")
	}

}

func TestMutationIterator(t *testing.T) {
	models := StrictAtom("model", true, nil)

	_, ok := models.Mutate("admin")
	if !ok {
		t.Fatal("Unable to mutate value")
	}

	_, ok = models.Mutate("users")
	if !ok {
		t.Fatal("Unable to mutate value")
	}

	_, ok = models.Mutate("groups")
	if !ok {
		t.Fatal("Unable to mutate value")
	}

	_, ok = models.Mutate(1)
	if ok {
		t.Fatal("Unable to mutate value")
	}

	at := models.Stamp().Add(-1 * time.Second)

	models.Mutate("laggies")

	itr := FiniteEventIterator(models, at)

	for itr.Next() == nil {
		t.Logf("Current State: %s stamp: %s againts %s", itr.Event(), itr.Event().Stamp(), at)
	}
}
