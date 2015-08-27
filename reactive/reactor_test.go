package reactive

import (
	"sync"
	"testing"

	"github.com/influx6/flux"
)

func TestBasicTimeBehaviour(t *testing.T) {
	nums := flux.ReactIdentity()
	memos := TimeTransform(nums)

	ws := new(sync.WaitGroup)

	ws.Add(2)

	memos.React(func(r flux.Reactor, err error, d interface{}) {
		if _, ok := d.(int); !ok {
			flux.FatalFailed(t, "Expected type int but got %+v", d)
		}
		ws.Done()
	}, true)

	nums.Send(20)
	flux.LogPassed(t, "Delivered 20 to num reactor")
	nums.Send(30)
	flux.LogPassed(t, "Delivered 30 to num reactor")

	ws.Wait()

	memos.Close()
	nums.Close()
}

func TestReactor(t *testing.T) {
	nums := flux.ReactIdentity()
	memos := TimeTransform(nums)

	ws := new(sync.WaitGroup)

	ws.Add(5)

	memos.React(func(r flux.Reactor, err error, d interface{}) {
		ws.Done()
	}, true)

	nums.Send(20)
	nums.Send(40)
	nums.Send(30)
	nums.Send(60)
	nums.Send(50)
	flux.LogPassed(t, "Delivered 4 numbers to num:reactor")

	ws.Wait()

	cop, err := memos.TimeLord()

	if err != nil {
		flux.FatalFailed(t, "Unable to create timecop: %+s", err)
	}

	flux.LogPassed(t, "Successfully created timelord for all changes")

	ws.Add(3)

	cop.RewindTime()
	cop.RewindTime()
	flux.LogPassed(t, "Successfully rewinded time to 2 past point")
	cop.ForwardTime()
	flux.LogPassed(t, "Successfully rewinded time to 1 future point")

	ws.Wait()

	memos.Close()
	nums.Close()
}
