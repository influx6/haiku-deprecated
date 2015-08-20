package flux

import (
	"sync"
	"testing"
	"time"
)

func TestWhileTicker(t *testing.T) {
	w := While(2 * time.Second)
	w.Start()

	ws := new(sync.WaitGroup)
	ws.Add(2)

	count := 0
	GoDefer("WaitTickerTest", func() {
		for w.Signals() != nil {
			select {
			case <-w.Signals():
				t.Logf("Received Signal %d", count)
				count++
				ws.Done()
			}
		}
		log.Info("closing loop")
	})

	ws.Wait()
	w.Stop()

}

func TestResetTimer(t *testing.T) {
	ws := new(sync.WaitGroup)
	rs := NewResetTimer(func() {
		t.Log("Initing reset")
	}, func() {
		t.Log("Finishing reset")
		ws.Done()
	}, time.Duration(600)*time.Millisecond, true, true)

	ws.Add(2)

	go func() {
		time.Sleep(time.Duration(650) * time.Millisecond)
		rs.Add()
	}()

	ws.Wait()
	rs.Close()
}

func TestOneTimeWaiter(t *testing.T) {
	ms := time.Duration(3) * time.Second

	w := NewTimeWait(0, ms)
	ws := new(sync.WaitGroup)

	t.Log("Before Add: Count:", w.Count())
	w.Add()
	ws.Add(1)
	t.Log("After 2 Add: Count:", w.Count())

	w.Then().When(func(v interface{}, _ ActionInterface) {
		_, ok := v.(int)

		if !ok {
			t.Fatal("Waiter completed with non-int value:", v)
		}

		t.Log("TimeWaiter Finished")
		ws.Done()
	})

	if w == nil {
		t.Fatal("Unable to create waiter")
	}

	if w.Count() < 0 {
		t.Fatal("Waiter completed before use")
	}

	if w.Count() < 1 {
		t.Fatalf("TimeWaiter count is at %d below 1 after two calls to Add()", w.Count())
	}

	if w.Count() < -1 {
		t.Fatal("Waiter just did a bad logic and went below -1")
	}

	// t.Log("Calling Done()")
	w.Done()

	ws.Wait()
}

func TestTimeWaiter(t *testing.T) {
	ms := time.Duration(1) * time.Second

	w := NewTimeWait(0, ms)
	ws := new(sync.WaitGroup)

	t.Log("Before Add: Count:", w.Count())
	ws.Add(1)

	w.Then().When(func(v interface{}, _ ActionInterface) {
		_, ok := v.(int)

		if !ok {
			t.Fatal("Waiter completed with non-int value:", v)
		}

		t.Log("TimeWaiter Finished")
		ws.Done()
	})

	if w == nil {
		t.Fatal("Unable to create waiter")
	}

	if w.Count() < 0 {
		t.Fatal("Waiter completed before use")
	}

	t.Log("Before Add: Count:", w.Count())
	w.Add()
	w.Add()
	t.Log("After 2 Add: Count:", w.Count())

	if w.Count() < 1 {
		t.Fatalf("TimeWaiter count is at %d below 1 after two calls to Add()", w.Count())
	}

	if w.Count() < -1 {
		t.Fatal("Waiter just did a bad logic and went below -1")
	}

	// w.Done()

	ws.Wait()
}

func TestWaiter(t *testing.T) {
	w := NewWait()

	w.Then().When(func(v interface{}, _ ActionInterface) {
		_, ok := v.(int)

		if !ok {
			t.Fatal("Waiter completed with non-int value:", v)
		}

		t.Log("SimpleWaiter Finished")
	})

	if w == nil {
		t.Fatal("Unable to create waiter")
	}

	if w.Count() < 0 {
		t.Fatal("Waiter completed before use")
	}

	w.Add()

	if w.Count() < 1 {
		t.Fatalf("Waiter count is still at %v despite call to Add()", w.Count())
	}

	cu := w.Count()
	t.Log("SimpleWait Count:", cu)

	w.Done()

	if w.Count() > 1 {
		t.Fatalf("Waiter count (%v) is still greater (%v) despite call to Done()", w.Count(), cu)
	}

	w.Done()
	w.Done()

	if w.Count() < -1 {
		t.Fatal("Waiter just did a bad logic and went below -1")
	}
}
