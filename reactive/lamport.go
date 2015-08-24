package reactive

import "time"

type (
	//Lamport provides a simple lamport timer which ensures the next time is greater than the previous time
	Lamport struct {
		seq              int
		lastSeenDuration time.Duration
		lastSeenTime     time.Time
		lastSeenSeq      int
		offset           time.Duration
	}

	//TimeStamp provide stamp details
	TimeStamp struct {
		Stamp time.Time
		Seq   int
	}

	//Timer provides an interface for time generators
	Timer interface {
		GetTime() *TimeStamp
		AdjustTime(time.Time)
	}
)

var (
	//EPOCH since unix-set date time
	EPOCH, _ = time.Parse("UTF", "Wed, 01 Jan 2014 00:00:00 GMT")
)

//NewLamport returns a lamport timer
func NewLamport(op time.Duration) *Lamport {
	return &Lamport{
		offset: op,
	}
}

//Minutes returns the minute since EPOCH with an offset
func (l *Lamport) minutes() time.Duration {
	ms := int64(time.Now().Sub(EPOCH).Minutes()) / 60000000
	ms += int64(l.offset.Seconds())
	return time.Duration(ms) * time.Minute
}

//Seconds returns the minute since EPOCH with an offset
func (l *Lamport) seconds() time.Duration {
	return time.Duration(l.minutes().Seconds())
}

//GetTime returns a new timestamp
func (l *Lamport) GetTime() *TimeStamp {
	min := l.minutes()

	if min > l.lastSeenDuration {
		l.lastSeenDuration = min
	}

	if min < l.lastSeenDuration {
		min = l.lastSeenDuration
	}

	var mst time.Time

	if l.lastSeenTime.IsZero() {
		mst = time.Now().Add(min)
	} else {
		mst = l.lastSeenTime.Add(min)
	}

	l.seq++
	l.lastSeenTime = mst

	return &TimeStamp{
		Stamp: mst,
		Seq:   l.seq,
	}

}

//AdjustTime adjust the lamport stamp offset by the total difference of the given time from the current time
func (l *Lamport) AdjustTime(ms time.Time) {
	jo := time.Duration(ms.Sub(time.Now()).Minutes()) * time.Minute
	l.offset = jo
}
