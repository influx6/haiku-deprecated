package reactive

import (
	"math/rand"
	"testing"
	"time"
)

func TestMutationRange(t *testing.T) {
	models := StrictAtom("model", true)
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
	models := StrictAtom("model", true)
	models.Mutate("admin")
	models.Mutate("users")
	grps, _ := models.Mutate("groups")
	models.Mutate("laggies")

	itr := NewIterator(NewMutationRange(models, grps))
	count := 4

	fw := 0
	for itr.Next() == nil {
		fw++
	}

	if fw != count {
		t.Fatal("Forward-Iterator release more than expected:", fw)
	}

	rti := itr.Reverse()
	fw = 0

	for rti.Next() == nil {
		fw++
	}

	if fw != count {
		t.Fatal("Reverse-Iterator release more than expected:", fw)
	}
}

func TestRestrictedListManager(t *testing.T) {

	lm := NewListManager(4, StrictAtom(20, true))

	if lm.MaxRange() > 4 {
		t.Fatal("Incorrect mutation range returned, expected 4")
	}

	_, err := lm.First()

	if err != nil {
		t.Fatal("Mutation should already have one mutation allocated")
	}

	_, err = lm.Last()

	if err != nil {
		t.Fatal("Mutation should already have one mutation allocated")
	}

	if lm.Size() <= 0 {
		t.Fatal("Manager should not be empty at this point")
	}

	_, ok := lm.Mutate(200)

	if !ok {
		t.Fatal("Unable to mutate list mutation value")
	}

	_, ok = lm.Mutate("400")

	if ok {
		t.Fatal("Bad: Allowed to mutate list mutation value from int to string")
	}

	fo, _ := lm.First()

	if fo.Value() != 20 {
		t.Fatal("Mutation has wrong first mutation value")
	}

	lo, _ := lm.Last()

	if lo.Value() != 200 {
		t.Fatal("Mutation has wrong last mutation value")
	}
}

func TestUnRestrictedListManager(t *testing.T) {

	lm := NewListManager(4, nil)

	if lm.MaxRange() > 4 {
		t.Fatal("Incorrect mutation range returned, expected 4")
	}

	if lm.Size() > 0 {
		t.Fatal("Manager should be empty at this point")
	}

	_, ok := lm.Mutate(20)

	if !ok {
		t.Fatal("Unable to mutate within list")
	}

	_, ok = lm.Mutate("words")

	if !ok {
		t.Fatal("Unable to mutate within list from int to string")
	}

	fo, _ := lm.First()

	if fo.Value() != 20 {
		t.Fatal("Incorrect first value,expected 20")
	}

	lo, _ := lm.Last()

	if lo.Value() != "words" {
		t.Fatal("Incorrect last value,expected 'words'")
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
		t.Fatal(err)
	}

	cur := 0
	for itr.Next() == nil {
		cur++
	}

	if cur > 5 {
		t.Fatal("Count is greater than 5")
	}
}

// func TestReactor(t *testing.T) {
// 	react := TypeReactor(100, "kind")
//
// 	if react == nil {
// 		t.Fatal("Unable to create reactor")
// 	}
//
// 	if m := react.Get(); m != "kind" {
// 		t.Fatal("Initial Value incorrect: ", m)
// 	}
//
// 	onems := time.Now()
//
// 	react.Set("love")
//
// 	if m := react.Get(); m != "love" {
// 		t.Fatal("Initial Value incorrect: ", m)
// 	}
//
// 	l, ok := react.Mutate("sucker")
//
// 	if m := l.Value(); m != "sucker" {
// 		t.Fatal("Initial Value incorrect: ", m, ok)
// 	}
//
// 	lastSave, err := react.SnapFrom(Event(onems.Unix()))
//
// 	if err != nil {
// 		t.Fatal("Unable to get history: ", err.Error())
// 	}
//
// 	log.Println(lastSave)
//
// }
