package rui

import (
	"log"
	"testing"
	"time"

	"github.com/influx6/flux"
	"github.com/influx6/prox/reactive"
)

type sampleRob struct {
	// DataTrees
	Name reactive.Observers
	Age  reactive.Observers
	Date time.Time
}

func TestBasicRendable(t *testing.T) {
	name, _ := reactive.ObserveTransform("Alex", false)
	age, _ := reactive.ObserveTransform(1, false)

	bob := &sampleRob{
		Name: name,
		Age:  age,
		Date: time.Now(),
	}

	_, err := StructTree(bob)

	if err != nil {
		flux.FatalFailed(t, "Unable to create struct tree: %s", err.Error())
	}

	bob.Name.Set("Joe")
}

func TestTemplateRendering(t *testing.T) {

	tmpl, err := SourceTemplator("base.tml", `
    <div>{{.Name}}</div>
    <div>{{.Age}}</div>
  `)

	if err != nil {
		flux.FatalFailed(t, "Unable to create template gen tree: %s", err.Error())
	}

	name, _ := reactive.ObserveTransform("Alex", false)
	age, _ := reactive.ObserveTransform(1, false)

	tol, err := tmpl.Build(&sampleRob{
		Name: name,
		Age:  age,
		Date: time.Now(),
	})

	if err != nil {
		flux.FatalFailed(t, "Unable to create templateRenderer: %s", err.Error())
	}

	log.Printf("Render: %s", tol.Render())

	tol.React(func(r flux.Reactor, err error, d interface{}) {
		log.Printf("tol: %s", d)
	}, true)

	name.Set("ron!")
}
