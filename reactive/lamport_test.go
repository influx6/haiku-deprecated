package reactive

import (
	"testing"
	"time"
)

func TestLamport(t *testing.T) {
	lm := NewLamport(0)

	ms := lm.GetTime()
	if ms.Stamp.After(lm.GetTime().Stamp) {
		t.Fatalf("%s is greater than the next time", ms)
	}

	//adjust times by a constant of 20min
	lm.AdjustTime(time.Now().Add(20 * time.Minute))

	ds := lm.GetTime().Stamp.Sub(lm.GetTime().Stamp)

	if ds == (20 * time.Minute) {
		t.Fatalf("%s difference shout be 20min apart", ds)
	}

}

func Test2Lamport(t *testing.T) {
	m1 := NewLamport(0)
	m2 := NewLamport(0)

	ms := m1.GetTime()
	md := m2.GetTime()

	if ms.Stamp.Equal(md.Stamp) {
		t.Fatalf("%s and %s should not be equal", ms, md)
	}

}
