package rui

import (
	"testing"

	"github.com/influx6/flux"
	"github.com/influx6/prox/reactive"
)

type foo struct {
	Name reactive.Observers
}

func TestDataTreeWithStruct(t *testing.T) {

	name, err := reactive.ObserveTransform("yo!", false)

	if err != nil {
		flux.FatalFailed(t, "Error creating observer: %s", err.Error())
	}

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

	dog, err := reactive.ObserveTransform("bullwilder", false)

	if err != nil {
		flux.FatalFailed(t, "Error creating observer: %s", err.Error())
	}

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

	dog, err := reactive.ObserveTransform("bullwilder", false)

	if err != nil {
		flux.FatalFailed(t, "Error creating observer: %s", err.Error())
	}

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
