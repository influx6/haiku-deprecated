package dom

import (
	"errors"
	"fmt"

	"github.com/go-humble/detect"
	"github.com/gopherjs/gopherjs/js"
	"github.com/influx6/flux"
	"github.com/influx6/haiku/views"
	"honnef.co/go/js/dom"
)

/*
  Path reactors dont try to be smart,they are simple, dum and do their best at what they do type of system, they simply report to you when certain change occurs with the system they are watching in the browser dom and api. They produce continous stream of events when changes occur
*/

//ErrNotSupported is returned when a feature requested is not supported by the environment
var ErrNotSupported = errors.New("Feature not supported")

var browserSupportsPushState = (js.Global.Get("onpopstate") != js.Undefined) &&
	(js.Global.Get("history") != js.Undefined) &&
	(js.Global.Get("history").Get("pushState") != js.Undefined)

// PathObserver represent any continouse changing route path by the browser
type PathObserver struct {
	flux.Reactor
}

// Path returns a new PathObserver instance
func Path() *PathObserver {
	return &PathObserver{
		Reactor: flux.ReactIdentity(),
	}
}

// Follow creates a Pathspec from the hash and path and sends it
func (p *PathObserver) Follow(hash, path string) {
	p.FollowSpec(&PathSpec{Hash: hash, Path: path})
}

// FollowSpec just passes down the Pathspec
func (p *PathObserver) FollowSpec(ps *PathSpec) {
	p.Send(ps)
}

// PathSpec represent the current path and hash values
type PathSpec struct {
	Hash string
	Path string
}

// String returns the hash and path
func (p *PathSpec) String() string {
	return fmt.Sprintf("%s%s", p.Path, p.Hash)
}

// HashChangePath returns a path observer path changes
func HashChangePath() *PathObserver {
	panicBrowserDetect()
	path := Path()

	js.Global.Set("onhashchange", func() {
		go func() {
			loc := js.Global.Get("location")
			pathn := loc.Get("pathname").String()
			hash := loc.Get("hash").String()
			path.Follow(hash, pathn)
		}()
	})

	return path
}

// PopStatePath returns a path observer path changes
func PopStatePath() (*PathObserver, error) {
	panicBrowserDetect()

	if !browserSupportsPushState {
		return nil, ErrNotSupported
	}

	path := Path()

	js.Global.Set("onpopstate", func() {
		go func() {
			loc := js.Global.Get("location")
			pathn := loc.Get("pathname").String()
			hash := loc.Get("hash").String()
			path.Follow(hash, pathn)
		}()
	})

	return path, nil
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

// Page provides the concrete provider for managing a whole website or View
// you dont need two,just one is enough to manage the total web view of your app / site
// It ties directly into the page hash or popstate location to provide consistent updates
type Page struct {
	*views.StateEngine
	path  *PathObserver
	views []*ViewComponent
}

// NewPage returns the new state engine powered page
func NewPage(p *PathObserver) *Page {
	pg := &Page{
		StateEngine: views.NewStateEngine(),
		path:        p,
		views:       make([]*ViewComponent, 0),
	}

	pg.All(".")
	return pg
}

// MyPage builds a new page with the appropriate path manager based on browser support, if the push and popstate were supported it will use that else use a hashpath manager
func MyPage() *Page {
	var p *PathObserver

	if po, err := PopStatePath(); err == nil {
		p = po
	} else {
		p = HashChangePath()
	}

	return NewPage(p)
}

// ErrBadSelector is used to indicate if the selector returned no result
var ErrBadSelector = errors.New("Selector returned nil")

// Mount adds a component into the page for handling/managing of visiblity and
// gets the dom referenced by the selector using QuerySelector and returns an error if selector gave no result
func (p *Page) Mount(selector, addr string, v views.Views) error {
	n := dom.GetWindow().Document().QuerySelector(selector)

	if n == nil {
		return ErrBadSelector
	}

	bv := BasicView(v)
	p.views = append(p.views, bv)
	p.AddView(addr, v)
	bv.Mount(n)
	return nil
}

// AddView adds a view to the page
func (p *Page) AddView(addr string, v views.Views) {
	p.UseState(addr, v)
}

// Address returns the current path and hash of the location api
func (p *Page) Address() (string, string) {
	loc := js.Global.Get("location")
	pathn := loc.Get("pathname").String()
	hash := loc.Get("hash").String()
	return pathn, hash
}
