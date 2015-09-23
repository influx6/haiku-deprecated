package reactive

import (
	"html/template"
	"testing"
	"time"

	"github.com/influx6/flux"
)

//basic model or template data that has reactive elements
type sampleRob struct {
	// DataTrees
	Name Observers
	Age  Observers
	Date time.Time
}

// TestBasicRendable creates a basic reactive struct whoes data contain some reactive elements which are then used to build a datatree which listens to each of these elements for change
func TestBasicRendable(t *testing.T) {
	name, _ := ObserveTransform("Alex", false)
	age, _ := ObserveTransform(1, false)

	bob := &sampleRob{
		Name: name,
		Age:  age,
		Date: time.Now(),
	}

	box, err := StructTree(bob)

	if err != nil {
		flux.FatalFailed(t, "Unable to create struct tree: %s", err.Error())
	}

	box.React(func(r flux.Reactor, err error, d interface{}) {
		if d != bob {
			flux.FatalFailed(t, "target and structree reply is not equal: %s %s", d, bob)
		}
	}, true)

	bob.Name.Set("Joe")
}

// TestTemplateRendering combines the template renderer for reactive structs with the data tree provider which is generated when a call to Build() with the supplied struct type (interface{})
func TestTemplateRendering(t *testing.T) {
	tml, err := template.New("base.tml").Parse(`
    <div>{{.Name}}</div>
    <div>{{.Age}}</div>
    <div>{{.Date}}</div>
	`)

	if err != nil {
		flux.FatalFailed(t, "Unable to create template gen tree: %s", err.Error())
	}

	name, _ := ObserveTransform("Alex", false)
	age, _ := ObserveTransform(1, false)

	tol, err := NewTempler(tml, &sampleRob{
		Name: name,
		Age:  age,
		Date: time.Now(),
	})

	if err != nil {
		flux.FatalFailed(t, "Unable to create templateRenderer: %s", err.Error())
	}

	cur := tol.Render()

	tol.React(func(r flux.Reactor, err error, d interface{}) {
		if cur == d {
			flux.FatalFailed(t, "both %s and %s should not be equal", cur, d)
		}

		flux.LogPassed(t, "both %s and %s are not be equal", cur, d)
	}, true)

	name.Set("Ron!")
	name.Set("Ron!")
}
