package rui

import (
	"testing"

	"github.com/influx6/flux"
	"github.com/influx6/prox/reactive"
)

type foo struct {
	Name reactive.Observers
}

func TestDataTreeReaction(t *testing.T) {

	name, _ := reactive.ObserveTransform("yo!", false)

	bar := &foo{name}

	tree := NewDataTree()

	if err := RegisterStructObservers(tree, bar); err != nil {
		flux.FatalFailed(t, "Unable to register observers: %s", err.Error())
	}

	tree.React(func(r flux.Reactor, err error, m interface{}) {
		dm, err := m.(DataTrees).Track("Name")

		if err != nil {
			flux.FatalFailed(t, "Error getting observer with field Name: %s", err.Error())
		}

		if dm.Get() != "dude!" {
			flux.FatalFailed(t, "data with incorrect value expected 'dude!': %s", dm.Get())
		}

		flux.LogPassed(t, "datatree has change with Name field with value: %s", dm.Get())
	}, true)

	bar.Name.Set("dude!")
}

func TestDataTreeWithStruct(t *testing.T) {

	name, _ := reactive.ObserveTransform("yo!", false)

	bar := &foo{name}

	tree := NewDataTree()

	if err := RegisterStructObservers(tree, bar); err != nil {
		flux.FatalFailed(t, "Unable to register observers: %s", err.Error())
	}

	flux.LogPassed(t, "Successfully completed RegisterStructObserver call")

	//do we a kind registered into the tree
	if ok := tree.Tracking("Name"); !ok {
		flux.FatalFailed(t, "Name is not in trackers")
	}

	flux.LogPassed(t, "Successfully tracking Name field")
}

func TestDataTreeWithMap(t *testing.T) {

	dog, _ := reactive.ObserveTransform("bullwilder", false)

	bar := map[string]interface{}{
		"jerry": dog,
		"tom":   "astracian",
	}

	tree := NewDataTree()

	if err := RegisterMapObservers(tree, bar); err != nil {
		flux.FatalFailed(t, "Unable to register observers: %s", err.Error())
	}

	flux.LogPassed(t, "Successfully completed RegisterStructObserver call")

	//do we a kind registered into the tree
	if ok := tree.Tracking("jerry"); !ok {
		flux.FatalFailed(t, "Jerry is not in trackers")
	}

	//we should not have tom in the tree
	if ok := tree.Tracking("tom"); ok {
		flux.FatalFailed(t, "tom is not a reactor")
	}

	flux.LogPassed(t, "Successfully tracking Jerry field")
}

func TestDataTreeWithList(t *testing.T) {

	dog, _ := reactive.ObserveTransform("bullwilder", false)

	bar := []interface{}{"astracian", dog}

	tree := NewDataTree()

	if err := RegisterListObservers(tree, bar); err != nil {
		flux.FatalFailed(t, "Unable to register observers: %s", err.Error())
	}

	flux.LogPassed(t, "Successfully completed RegisterStructObserver call")

	//do we a kind registered into the tree
	if ok := tree.Tracking("1"); !ok {
		flux.FatalFailed(t, "1 is not in trackers")
	}

	//we should not have tom in the tree
	if ok := tree.Tracking("0"); ok {
		flux.FatalFailed(t, "0 is not a reactor")
	}

	flux.LogPassed(t, "Successfully tracking 0 index field")
}
