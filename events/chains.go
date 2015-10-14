package events

import "sync"

// hodom "honnef.co/go/js/dom"

// NextHandler provides next call for flat chains
type NextHandler func(Event)

// FlatHandler provides a handler for flatchain
type FlatHandler func(Event, NextHandler)

//FlatChains define a simple flat chain
type FlatChains interface {
	HandleContext(Event)
	Next(FlatHandler) FlatChains
	Chain(FlatChains) FlatChains
	NChain(FlatChains) FlatChains
	Bind(rnx NextHandler) FlatChains
	useNext(FlatChains)
	usePrev(FlatChains)
	UnChain()
}

// FlatChain provides a simple middleware like
type FlatChain struct {
	op         FlatHandler
	prev, next FlatChains
}

//FlatChainIdentity returns a chain that calls the next automatically
func FlatChainIdentity() FlatChains {
	return NewFlatChain(func(c Event, nx NextHandler) {
		nx(c)
	})
}

//NewFlatChain returns a new flatchain instance
func NewFlatChain(fx FlatHandler) *FlatChain {
	return &FlatChain{
		op: fx,
	}
}

// UnChain unlinks the current chain from the set and reconnects the others
func (r *FlatChain) UnChain() {
	prev := r.prev
	next := r.next

	if prev != nil && next != nil {
		prev.useNext(next)
		next.usePrev(prev)
		return
	}

	prev.useNext(nil)
}

// Bind provides a wrapper over the Next binder function call
func (r *FlatChain) Bind(rnx NextHandler) FlatChains {
	return r.Next(func(ev Event, nx NextHandler) {
		rnx(ev)
		nx(ev)
	})
}

// Next allows the chain of the function as a FlatHandler
func (r *FlatChain) Next(rnx FlatHandler) FlatChains {
	nx := NewFlatChain(rnx)
	return r.NChain(nx)
}

// Chain sets the next flat chains else passes it down to the last chain to set as next chain,returning itself
func (r *FlatChain) Chain(rx FlatChains) FlatChains {
	if r.next == nil {
		rx.usePrev(r)
		r.useNext(rx)
	} else {
		r.next.Chain(rx)
	}
	return r
}

// NChain sets the next flat chains else passes it down to the last chain to set as next chain,returning the the supplied chain
func (r *FlatChain) NChain(rx FlatChains) FlatChains {
	if r.next == nil {
		r.useNext(rx)
		return rx
	}

	return r.next.NChain(rx)
}

// HandleContext calls the next chain if any
func (r *FlatChain) HandleContext(c Event) {
	r.op(c, func(c Event) {
		if r.next != nil {
			r.next.HandleContext(c)
		}
	})
}

// useNext swaps the next chain with the supplied
func (r *FlatChain) useNext(fl FlatChains) {
	r.next = fl
}

// usePrev swaps the previous chain with the supplied
func (r *FlatChain) usePrev(fl FlatChains) {
	r.prev = fl
}

//ChainFlats chains second flats to the first flatchain and returns the first flatchain
func ChainFlats(mo FlatChains, so ...FlatChains) FlatChains {
	for _, sof := range so {
		func(do FlatChains) {
			mo.Chain(do)
		}(sof)
	}
	return mo
}

//ElemEventMux represents a stanard callback function for dom events
type ElemEventMux func(Event, func())

//ListenerStack provides addition of functions into a stack
type ListenerStack struct {
	listeners []ElemEventMux
	lock      sync.RWMutex
}

//NewListenerStack returns a new ListenerStack instance
func NewListenerStack() *ListenerStack {
	return &ListenerStack{
		listeners: make([]ElemEventMux, 0),
	}
}

//Clear flushes the stack listener
func (f *ListenerStack) Clear() {
	f.lock.Lock()
	f.listeners = f.listeners[0:0]
	f.lock.Unlock()
}

//Size returns the total number of listeners
func (f *ListenerStack) Size() int {
	f.lock.RLock()
	sz := len(f.listeners)
	f.lock.RUnlock()
	return sz
}

//Add adds a function into the stack
func (f *ListenerStack) Add(fx ElemEventMux) int {
	var ind int

	f.lock.RLock()
	ind = len(f.listeners)
	f.lock.RUnlock()

	f.lock.Lock()
	f.listeners = append(f.listeners, fx)
	f.lock.Unlock()

	return ind
}

// DeleteIndex removes the function at the provided index
func (f *ListenerStack) DeleteIndex(ind int) {

	if ind <= 0 && len(f.listeners) <= 0 {
		return
	}

	f.lock.Lock()
	copy(f.listeners[ind:], f.listeners[ind+1:])
	f.lock.Unlock()

	f.lock.RLock()
	f.listeners[len(f.listeners)-1] = nil
	f.lock.RUnlock()

	f.lock.Lock()
	f.listeners = f.listeners[:len(f.listeners)-1]
	f.lock.Unlock()
}

//Each runs through the function lists and executing with args
func (f *ListenerStack) Each(d Event) {
	if f.Size() <= 0 {
		return
	}

	f.lock.RLock()

	// var ro sync.Mutex
	var stop bool

	for _, fx := range f.listeners {
		if stop {
			break
		}
		//TODO: is this critical that we send it into a goroutine with a mutex?
		fx(d, func() {
			// ro.Lock()
			stop = true
			// ro.Unlock()
		})
	}

	f.lock.RUnlock()
}
