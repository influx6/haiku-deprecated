package reactive

import (
	"log"
	"testing"
	"time"
)

func TestImmutable(t *testing.T) {
	models, err := Transform("model", false)

	if err != nil {
		t.Fatal(err)
	}

	if models.Get() != "model" {
		t.Fatal("Wrong returned value:", models.Get())
	}

	nc := models.Channel()

	models.Set("user")

	if data := <-nc; "user" != data {
		t.Fatal("Wrong channel returned value:", data)
	}

	if models.Get() == "model" {
		t.Fatal("Wrong returned value:", models.Get())
	}

}

func TestReactor(t *testing.T) {
	react := TypeReactor(100, "kind")

	if react == nil {
		t.Fatal("Unable to create reactor")
	}

	if m := react.Get(); m != "kind" {
		t.Fatal("Initial Value incorrect: ", m)
	}

	onems := time.Now()

	react.Set("love")

	if m := react.Get(); m != "love" {
		t.Fatal("Initial Value incorrect: ", m)
	}

	l, ok := react.Mutate("sucker")

	if m := l.Value(); m != "sucker" {
		t.Fatal("Initial Value incorrect: ", m, ok)
	}

	lastSave, err := react.SnapFrom(Event(onems.Unix()))

	if err != nil {
		t.Fatal("Unable to get history: ", err.Error())
	}

	log.Println(lastSave)

}
