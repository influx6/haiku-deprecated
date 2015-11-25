package views

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-humble/detect"
	"github.com/gopherjs/gopherjs/js"
	"github.com/influx6/haiku/jsutils"
	"github.com/influx6/haiku/pub"
)

/*
  Path reactors dont try to be smart,they are simple, dum and do their best at what they do type of system, they simply report to you when certain change occurs with the system they are watching in the browser dom and api. They produce continous stream of events when changes occur
*/

// PathSpec represent the current path and hash values
type PathSpec struct {
	Hash     string
	Path     string
	Sequence string
}

// String returns the hash and path
func (p *PathSpec) String() string {
	return fmt.Sprintf("%s%s", p.Path, p.Hash)
}

// PathObserver represent any continouse changing route path by the browser
type PathObserver struct {
	pub.Publisher
	usingHash bool
}

// Path returns a new PathObserver instance
func Path() *PathObserver {
	return &PathObserver{
		Publisher: pub.Identity(),
	}
}

// Follow creates a Pathspec from the hash and path and sends it
func (p *PathObserver) Follow(path, hash string) {
	cleanHash := strings.Replace(hash, "#", ".", -1)
	cleanHash = strings.Replace(cleanHash, "/", ".", -1)
	p.FollowSpec(PathSpec{Hash: hash, Path: path, Sequence: cleanHash})
}

// FollowSpec just passes down the Pathspec
func (p *PathObserver) FollowSpec(ps PathSpec) {
	p.Send(ps)
}

// NotifyPage is used to notify a page of
func (p *PathObserver) NotifyPage(pg *Pages) {
	p.React(func(r pub.Publisher, _ error, d interface{}) {
		// if err != nil { r.SendError(err) }
		// log.Printf("will Sequence: %s", d)
		if ps, ok := d.(PathSpec); ok {
			// log.Printf("Sequence: %s", ps.Sequence)
			pg.All(ps.Sequence)
		}
	}, true)
}

// NotifyPartialPage is used to notify a Page using the page engine's Partial() activator
func (p *PathObserver) NotifyPartialPage(pg *Pages) {
	p.React(func(r pub.Publisher, _ error, d interface{}) {
		// if err != nil { r.SendError(err) }
		if ps, ok := d.(PathSpec); ok {
			pg.Partial(ps.Sequence)
		}
	}, true)
}

// HashChangePath returns a path observer path changes
func HashChangePath() *PathObserver {
	panicBrowserDetect()
	path := Path()
	path.usingHash = true

	js.Global.Set("onhashchange", func() {
		path.Follow(GetLocation())
	})

	return path
}

//ErrNotSupported is returned when a feature requested is not supported by the environment
var ErrNotSupported = errors.New("Feature not supported")

// PopStatePath returns a path observer path changes
func PopStatePath() (*PathObserver, error) {
	panicBrowserDetect()

	if !BrowserSupportsPushState() {
		return nil, ErrNotSupported
	}

	path := Path()

	js.Global.Set("onpopstate", func() {
		path.Follow(GetLocation())
	})

	return path, nil
}

// HistoryProvider wraps the PathObserver with methods that allow easy control of
// client location
type HistoryProvider struct {
	*PathObserver
}

// History returns a new PathObserver and depending on browser support will either use the
// popState or HashChange
func History() *HistoryProvider {
	pop, err := PopStatePath()

	if err != nil {
		pop = HashChangePath()
	}

	return &HistoryProvider{pop}
}

// Go changes the path of the current browser location depending on wether its underline
// observe is hashed based or pushState based,it will use SetDOMHash or PushDOMState appropriately
func (h *HistoryProvider) Go(path string) {
	if h.usingHash {
		SetDOMHash(path)
		return
	}
	PushDOMState(path)
}

// ErrBadSelector is used to indicate if the selector returned no result
var ErrBadSelector = errors.New("Selector returned nil")

// Pages provides the concrete provider for managing a whole website or View
// you dont need two,just one is enough to manage the total web view of your app / site
// It ties directly into the page hash or popstate location to provide consistent updates
type Pages struct {
	*StateEngine
	*HistoryProvider
	// views []Views
}

// Page returns the new state engine powered page
func Page() *Pages {
	return NewPage(History())
}

// NewPage returns the new state engine powered page
func NewPage(p *HistoryProvider) *Pages {
	pg := &Pages{
		StateEngine:     NewStateEngine(),
		HistoryProvider: p,
	}

	p.NotifyPage(pg)
	pg.All(".")
	return pg
}

// Mount adds a component into the page for handling/managing of visiblity and
// gets the dom referenced by the selector using QuerySelector and returns an error if selector gave no result
func (p *Pages) Mount(selector, addr string, v Views) error {
	// n := dom.GetWindow().Document().QuerySelector(selector)
	n := jsutils.GetDocument().Call("querySelector", selector)

	// log.Printf("we failed: %+s", n)

	if n == nil || n == js.Undefined {
		return ErrBadSelector
	}

	// log.Printf("no we did failed: %+s", n)

	// p.views = append(p.views, v)
	p.AddView(addr, v)
	v.Mount(n)
	return nil
}

// AddView adds a view to the page
func (p *Pages) AddView(addr string, v Views) {
	p.UseState(addr, v)
}

// Address returns the current path and hash of the location api
func (p *Pages) Address() (string, string) {
	return GetLocation()
}

// GetLocation returns the path and hash of the browsers location api else panics if not in a browser
func GetLocation() (string, string) {
	panicBrowserDetect()
	loc := js.Global.Get("location")
	path := loc.Get("pathname").String()
	hash := loc.Get("hash").String()
	return path, hash
}

// PushDOMState adds a new state the dom push history
func PushDOMState(path string) {
	panicBrowserDetect()
	js.Global.Get("history").Call("pushState", nil, "", path)
}

// SetDOMHash sets the dom location hash
func SetDOMHash(hash string) {
	panicBrowserDetect()
	js.Global.Get("location").Set("hash", hash)
}

func panicBrowserDetect() {
	if !detect.IsBrowser() {
		panic("expected to be used in a dom/browser env")
	}
}

// BrowserSupportsPushState checks if browser supports pushState
func BrowserSupportsPushState() bool {
	if !detect.IsBrowser() {
		return false
	}

	return (js.Global.Get("onpopstate") != js.Undefined) &&
		(js.Global.Get("history") != js.Undefined) &&
		(js.Global.Get("history").Get("pushState") != js.Undefined)
}
