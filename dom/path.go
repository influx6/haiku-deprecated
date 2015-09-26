package dom

import (
	"errors"
	"fmt"

	"github.com/go-humble/detect"
	"github.com/gopherjs/gopherjs/js"
	"github.com/influx6/flux"
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

// HashPath represents a path using url hash
type HashPath string

// HashChangePath returns a path observer path changes
func HashChangePath() *PathObserver {
	panicBrowserDetect()
	path := Path()

	js.Global.Set("onhashchange", func() {
		go func() {
			loc := js.Global.Get("location")
			pathn := loc.Get("pathname").String()
			hash := loc.Get("hash").String()
			path.Send(HashPath(fmt.Sprintf("%s%s", pathn, hash)))
		}()
	})

	return path
}

// PopPath represents a path using history popstate
type PopPath string

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
			path.Send(PopPath(fmt.Sprintf("%s%s", pathn, hash)))
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

// func init(){
//   if detect.IsBrowser() {
//
//   }
// }
