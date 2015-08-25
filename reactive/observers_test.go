package reactive

import (
	"testing"

	"github.com/influx6/flux"
)

func TestImmutable(t *testing.T) {
	models, err := ObserveTransform("model", false)

	if err != nil {
		t.Fatal(err)
	}

	if models.Get() != "model" {
		t.Fatal("Wrong returned value:", models.Get())
	}

	stream := flux.SignalCollector()
	channel := models.React(flux.ChannelReactProcessor(stream))

	models.Set("user")

	if data := <-stream.Signals(); "user" != data {
		t.Fatal("Wrong channel returned value:", data)
	}

	if models.Get() == "model" {
		t.Fatal("Wrong returned value:", models.Get())
	}

	channel.SendClose(0)

}
