package main

import (
	"fmt"
	"time"

	"github.com/influx6/flux"
	hdom "github.com/influx6/haiku/dom"
	dom "honnef.co/go/js/dom"
)

func buildDelagationTest(win dom.Window, doc dom.HTMLDocument) {
	container := doc.QuerySelector(".container")
	speak := container.QuerySelector(".speak p")

	//create our container for delegated events and extra juices
	elem := hdom.NewElement(container)

	clickEvent := hdom.NewElemEvent("click", "a.linkroller")
	//tell all events of this type to prevent default when they match
	clickEvent.PreventDefault = true

	//did we add the new event successfully?
	if ok := elem.WatchEvent(clickEvent); !ok {
		panic("failed to setup event")
	}

	clickEvent.Next(func(ev dom.Event, next hdom.NextHandler) {
		speak.SetInnerHTML(fmt.Sprintf("Event from %s target occured with live/delegation from %s", ev.Target(), container))
		win.Location().Hash = "#jqueried"
		go func() {
			<-time.After(800 * time.Millisecond)
			// speak.SetInnerHTML("")
			//set the html of the paragraph as empty
			elem.Html("", ".speak p")
			win.Location().Hash = ""
		}()
		//pass it down to the next dude waiting
		next(ev)
	})
}

func buildHistoryPathTest(win dom.Window, doc dom.HTMLDocument) {
	pop := doc.QuerySelector(".history .hash")
	pushhistory := hdom.HashChangePath()

	pushhistory.React(func(r flux.Reactor, err error, data interface{}) {
		pop.SetInnerHTML(fmt.Sprintf("%s", data))
	}, true)
}

func buildHistoryPushTest(win dom.Window, doc dom.HTMLDocument) {
	pop := doc.QuerySelector(".history .pushpop")
	pushhistory, err := hdom.PopStatePath()

	if err != nil {
		pop.SetInnerHTML("PopState not supported!")
		return
	}

	pushhistory.React(func(r flux.Reactor, err error, data interface{}) {
		pop.SetInnerHTML(fmt.Sprintf("%s", data))
	}, true)
}

//go:generate gopherjs build main.go -o ./app.js

func main() {

	win, ok := dom.GetWindow().(dom.Window)

	if !ok {

		panic("not in a browser-window")
	}

	doc, ok := win.Document().(dom.HTMLDocument)

	if !ok {
		panic("not in a browser-dom")
	}

	buildDelagationTest(win, doc)
	buildHistoryPathTest(win, doc)
	buildHistoryPushTest(win, doc)
}
