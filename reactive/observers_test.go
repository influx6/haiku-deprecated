package reactive

import (
	"sync"
	"testing"

	"github.com/influx6/flux"
)

func TestImmutable(t *testing.T) {
	var ws sync.WaitGroup
	ws.Add(1)

	models, err := ObserveTransform("model", false)

	if err != nil {
		t.Fatal(err)
	}

	if models.Get() != "model" {
		t.Fatal("Wrong returned value:", models.Get())
	}

	models.React(func(r flux.Reactor, err error, data interface{}) {
		if "user" != data {
			flux.FatalFailed(t, "Wrong channel returned value: %s", data)
		}
		ws.Done()
	}, true)

	models.Set("user")

	ws.Wait()

	if models.Get() == "model" {
		flux.FatalFailed(t, "Wrong returned value:", models.Get())
	}

	flux.LogPassed(t, "Completed atom updated")

	models.Close()

}
