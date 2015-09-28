package main

import (
	"fmt"
	"time"

	hdom "github.com/influx6/haiku/dom"
	dom "honnef.co/go/js/dom"
)

//go:generate gopherjs build main.go -o ./app.js

func main() {

	doc, ok := dom.GetWindow().Document().(dom.HTMLDocument)

	if !ok {
		panic("not in a browser-dom")
	}

	container := doc.QuerySelector(".container")
	speak := container.QuerySelector(".speak p")

	//create our container for delegated events and extra juices
	elem := hdom.NewElement(container)

	clickEvent := hdom.NewElemEvent("click", "a.linkroller")
	//tell all events of this type to prevent default when they match
	clickEvent.PreventDefault = true

	//did we add the new event successfully?
	if ok = elem.WatchEvent(clickEvent); !ok {
		panic("failed to setup event")
	}

	clickEvent.Next(func(ev dom.Event, next hdom.NextHandler) {
		speak.SetInnerHTML(fmt.Sprintf("Event from %s target occured with live/delegation from %s", ev.Target(), container))
		go func() {
			<-time.After(500 * time.Millisecond)
			// speak.SetInnerHTML("")
			//set the html of the paragraph as empty
			elem.Html("", ".speak p")
		}()
		//pass it down to the next dude waiting
		next(ev)
	})

}
