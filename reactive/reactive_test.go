package reactive

import "testing"

func TestImmutable(t *testing.T) {
	models, err := Transform("model")

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
