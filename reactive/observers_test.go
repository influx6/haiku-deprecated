package reactive

import (
	"encoding/json"
	"sync"
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/influx6/flux"
)

//foobar provides a basic struct we want to load yaml into
type foobar struct {
	Name *Observer `yaml:"name" json:"name"`
}

// TestObserverPayload test the capability of loading data into an Observer using yaml data payload
func TestObserverPayload(t *testing.T) {

	name, err := ObserveTransform("", false)

	if err != nil {
		flux.FatalFailed(t, "Unable to create observer: %s", err.Error())
	}

	bar := &foobar{Name: name}

	//load up a custom data in yaml format
	if err := yaml.Unmarshal([]byte(`name: "alex"`), bar); err != nil {
		flux.FatalFailed(t, "Unable to load data into observer using yaml: %s", err.Error())
	}

	flux.LogPassed(t, "Successfully executed yaml payload without error:", bar)

	//check if bar has the correct value
	if bar.Name.Get() != "alex" {
		flux.FatalFailed(t, "Observer contains wrong value; expected 'alex' got %s", bar.Name)
	}

	flux.LogPassed(t, "Successfully loaded and found correcte value: %s", bar.Name)

	//load up a custom data in json format
	if err := json.Unmarshal([]byte(`{"name": "john"}`), bar); err != nil {
		flux.FatalFailed(t, "Unable to load data into observer using json: %s", err.Error())
	}

	flux.LogPassed(t, "Successfully executed json payload without error:", bar)

	//check if bar has the correct value
	if bar.Name.Get() != "john" {
		flux.FatalFailed(t, "Observer contains wrong value; expected 'alex' got %s", bar.Name)
	}

	flux.LogPassed(t, "Successfully loaded and found correcte value: %s", bar.Name)

	// j, e := json.Marshal(bar)
	// log.Printf("json: %s %s", j, e)
	// y, e := yaml.Marshal(bar)
	// log.Printf("yaml: %s %s", y, e)
}

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
