package flux

import (
	"testing"
)

func TestApply(t *testing.T) {
	sc := NewStack(func(data interface{}, _ Stacks) interface{} {
		return data.(int) * 20
	})

	gs := sc.Stack(func(data interface{}, _ Stacks) interface{} {
		if data != 400 {
			t.Fatal("Value is incorrect:", data)
		}
		return data.(int) / 20
	}, true)

	_ = gs
	if val := sc.Apply(20); val != 20 {
		log.Debug("Return value is incorrect:", val)
	}
}

func TestLift(t *testing.T) {
	sc := NewStack(func(data interface{}, _ Stacks) interface{} {
		return data.(int) * 20
	})

	gs := sc.Stack(func(data interface{}, _ Stacks) interface{} {
		if data != 400 {
			t.Fatal("Value is incorrect:", data)
		}
		return data.(int) / 20
	}, true)

	if val := gs.Lift(20); val != 20 {
		log.Debug("lift: Return value is incorrect:", val)
	}
}

func TestLevitate(t *testing.T) {
	sc := NewStack(func(data interface{}, _ Stacks) interface{} {
		if data != 1 {
			t.Fatal("Value is incorrect,expected 1", data)
		}
		return data.(int) * 20
	})

	gs := sc.Stack(func(data interface{}, _ Stacks) interface{} {
		if data != 20 {
			t.Fatal("Value is incorrect:", data)
		}
		return data.(int) / 20
	}, true)

	if val := gs.Levitate(20); val != 1 {
		log.Debug("liftapply: Return value is incorrect:", val)
	}
}

func TestLiftApply(t *testing.T) {
	sc := NewStack(func(data interface{}, _ Stacks) interface{} {
		if data != 1 {
			t.Fatal("Value is incorrect,expected 1", data)
		}
		return data.(int) * 10
	})

	gs := sc.Stack(func(data interface{}, _ Stacks) interface{} {
		if data != 20 {
			t.Fatal("Value is incorrect:", data)
		}
		return data.(int) / 20
	}, true)

	if val := gs.LiftApply(20); val != 10 {
		log.Debug("liftapply: Return value is incorrect:", val)
	}
}

func TestCall(t *testing.T) {
	sc := NewStack(func(data interface{}, _ Stacks) interface{} {
		return data.(int) * 20
	})

	gs := sc.Stack(func(data interface{}, _ Stacks) interface{} {
		if data != 400 {
			t.Fatal("Value is incorrect:", data)
		}
		return data.(int) / 10
	}, true)

	_ = gs.Stack(func(data interface{}, _ Stacks) interface{} {
		if data != 40 {
			t.Fatal("Value is incorrect:", data)
		}
		return data.(int) / 20
	}, true)

	if val := sc.Call(20); val != 400 {
		log.Debug("Return value is incorrect:", val)
	}
}

func TestIsolate(t *testing.T) {
	sc := NewStack(func(data interface{}, _ Stacks) interface{} {
		return data.(int) * 20
	})

	_ = sc.Stack(func(data interface{}, _ Stacks) interface{} {
		t.Fatal("Stack was called,this is bad")
		return data.(int) / 20
	}, true)

	if val := sc.Isolate(20); val != 400 {
		log.Debug("Return value is incorrect:", val)
	}
}

func TestIdentity(t *testing.T) {
	gs := NewStack(func(data interface{}, _ Stacks) interface{} {
		return data.(int) * 20
	})

	if val := gs.Identity(20); val != 20 {
		log.Debug("Return value is incorrect:", val)
	}
}

func TestStackers(t *testing.T) {
	sc := NewStack(func(data interface{}, _ Stacks) interface{} {
		nm, ok := data.(int)
		if !ok {
			log.Debug("invalid type: expected int got %+s", data)
			return nil
		}

		return nm
	})

	g := sc.Stack(func(data interface{}, _ Stacks) interface{} {
		return data
	}, true)

	defer sc.Call(7)
	defer g.Call(82)
	// defer g.Close()

	_ = LogStack(g)

	_ = sc.Stack(func(data interface{}, _ Stacks) interface{} {
		return 20
	}, true)

	xres := sc.Call(1)
	yres := g.Call(30)

	if xres == yres {
		log.Debug("Equal unexpect values %d %d", xres, yres)
	}

	xres = sc.Call(40)
	g.Lift(20)

	if xres == yres {
		log.Debug("Equal unexpect values %d %d", xres, yres)
	}

	sc.Close()
}
