package views

import (
	"testing"

	"github.com/influx6/flux"
)

func TestStateEngineAll(t *testing.T) {
	var engine = NewStateEngine()

	// home :=
	home := engine.AddState("home")

	home.UseActivator(func() {
		flux.LogPassed(t, "Sucessfully activated Home")
	})

	home.Engine().AddState(".").UseActivator(func() {
		flux.LogPassed(t, "Sucessfully activated border")
	})

	home.Engine().AddState("swatch").UseActivator(func() {
		flux.LogPassed(t, "Sucessfully activated swatch")
	})

	err := engine.All(".home.swatch")

	if err != nil {
		flux.FatalFailed(t, "Unable to run full state: %s", err)
	}

}

func TestStateEnginePartial(t *testing.T) {
	var engine = NewStateEngine()

	home := engine.AddState("home")

	home.UseActivator(func() {
		flux.FatalFailed(t, "Should not have activated home")
	})

	home.Engine().AddState(".").UseActivator(func() {
		flux.FatalFailed(t, "Should not have activated border")
	})

	home.Engine().AddState("swatch").UseActivator(func() {
		flux.LogPassed(t, "Sucessfully activated swatch")
	})

	err := engine.Partial(".home.swatch")

	if err != nil {
		flux.FatalFailed(t, "Unable to run partial state: %s", err)
	}

}

func TestStateEngineDeactivate(t *testing.T) {
	var engine = NewStateEngine()

	home := engine.AddState("home")

	home.UseActivator(func() {
		flux.LogPassed(t, "Sucessfully activated home")
	})

	home.Engine().AddState("swatch").UseActivator(func() {
		flux.LogPassed(t, "Sucessfully activated swatch")
	}).UseDeactivator(func() {
		flux.LogPassed(t, "Sucessfully deactivated swatch")
	})

	err := engine.All(".home.swatch")

	if err != nil {
		flux.FatalFailed(t, "Unable to run full state: %s", err)
	}

	err = engine.All(".home")

	if err != nil {
		flux.FatalFailed(t, "Unable to run deactivate state: %s", err)
	}

}

func TestStateEngineRoot(t *testing.T) {
	var engine = NewStateEngine()

	home := engine.AddState(".")

	home.UseActivator(func() {
		flux.LogPassed(t, "Sucessfully activated home")
	})

	home.Engine().AddState(".").UseActivator(func() {
		flux.LogPassed(t, "Sucessfully activated swatch")
	}).UseDeactivator(func() {
		flux.LogPassed(t, "Sucessfully deactivated swatch")
	})

	err := engine.All(".")

	if err != nil {
		flux.FatalFailed(t, "Unable to run full state: %s", err)
	}
}
