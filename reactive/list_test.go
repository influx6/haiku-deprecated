package reactive

import (
	"math/rand"
	"testing"
	"time"

	"github.com/influx6/flux"
)

func TestMutationRange(t *testing.T) {
	models := Atom("model", true)
	models.Mutate("admin")
	models.Mutate("users")
	grps, _ := models.Mutate("groups")
	models.Mutate("laggies")

	mr := NewMutationRange(models, grps)

	if mr.Root() != models {
		t.Fatalf("MutationRange root incorrect: excpted %s got %s", models, mr.Root())
	}

	if mr.Tail() != grps {
		t.Fatalf("MutationRange root incorrect: excpted %s got %s", models, mr.Tail())
	}
}

func TestMutationIterator(t *testing.T) {
	models := Atom("model", true)
	models.Mutate("admin")
	models.Mutate("users")
	grps, _ := models.Mutate("groups")
	models.Mutate("laggies")

	itr := NewIterator(NewMutationRange(models, grps))
	count := 2

	fw := 0
	for itr.Next() == nil {
		if fw >= 2 {
			break
		}
		fw++
	}

	if fw != count {
		flux.FatalFailed(t, "Foward-Iterator release more than expected: %d", fw)
	}

	flux.LogPassed(t, "Iterator-Forward Works")

	// itr.Reset()
	itr.Reverse()

	fw = 0

	for itr.Next() == nil {
		fw++
	}

	if fw != count {
		flux.FatalFailed(t, "Reverse-Iterator release more than expected: %d", fw)
	}

	flux.LogPassed(t, "Iterator-Reverse Works")
}

func TestRestrictedListManager(t *testing.T) {

	lm := NewListManager(4, Atom(20, true))

	if lm.MaxRange() > 4 {
		flux.FatalFailed(t, "Incorrect mutation range returned, expected 4")
	}

	flux.LogPassed(t, "Correct mutation range returned, expected 4")

	_, err := lm.First()

	if err != nil {
		flux.FatalFailed(t, "Mutation first() failed, should already have one mutation allocated")
	}

	flux.LogPassed(t, "Mutation First Item Retrieved Successfully")

	_, err = lm.Last()

	if err != nil {
		flux.FatalFailed(t, "Mutation Last() call failed")
	}

	flux.LogPassed(t, "Mutation Last Item Retrieved Successfully")

	if lm.Size() <= 0 {
		flux.FatalFailed(t, "Manager should not be empty at this point")
	}

	flux.LogPassed(t, "Manager should has items,stated passed well!")

	_, ok := lm.Mutate(200)

	if !ok {
		flux.FatalFailed(t, "Unable to mutate list mutation value")
	}

	flux.LogPassed(t, "Mutate passed, new value is 200")

	_, ok = lm.Mutate("400")

	if ok {
		flux.FatalFailed(t, "Bad: Allowed to mutate list mutation value from int to string")
	}

	flux.LogPassed(t, "Mutate not allowed pased based on type, value is still 200")

	fo, _ := lm.First()

	if fo.Value() != 20 {
		flux.FatalFailed(t, "Mutation has wrong first mutation value")
	}
	flux.LogPassed(t, "Passed value check of first() in mutation list with value: %d", fo.Value())

	lo, _ := lm.Last()

	if lo.Value() != 200 {
		flux.FatalFailed(t, "Mutation has wrong last mutation value")
	}

	flux.LogPassed(t, "Passed value check of first() in mutation list with value: %d", fo.Value())
}

func TestUnRestrictedListManager(t *testing.T) {

	lm := NewListManager(4, nil)

	if lm.MaxRange() > 4 {
		flux.FatalFailed(t, "Incorrect mutation range returned, expected 4")
	}

	if lm.Size() > 0 {
		flux.FatalFailed(t, "Manager should be empty at this point")
	}

	_, ok := lm.Mutate(20)

	if !ok {
		flux.FatalFailed(t, "Unable to mutate within list")
	}

	_, ok = lm.Mutate("words")

	if !ok {
		flux.FatalFailed(t, "Unable to mutate within list from int to string")
	}

	fo, _ := lm.First()

	if fo.Value() != 20 {
		flux.FatalFailed(t, "Incorrect first value,expected 20")
	}

	lo, _ := lm.Last()

	if lo.Value() != "words" {
		flux.FatalFailed(t, "Incorrect last value,expected 'words'")
	}
}

func TestListManagerRetrievingSystem(t *testing.T) {
	lm := NewListManager(4, nil)

	count := 10

	for {
		if count <= 0 {
			break
		}

		lm.Mutate(count * rand.Intn(400))
		time.Sleep(300 * time.Millisecond)
		count--
	}

	itr, err := lm.SnapFrom(6 * time.Minute)

	if err != nil {
		flux.FatalFailed(t, "Error Occured snaping event timeline: %+s", err)
	}

	cur := 0
	for itr.Next() == nil {
		cur++
	}

	if cur > 5 {
		flux.FatalFailed(t, "Count is greater than 5")
	}
}
