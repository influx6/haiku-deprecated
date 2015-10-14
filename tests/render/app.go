package main

import (
	"log"
	"time"

	"github.com/influx6/haiku/trees"
	"github.com/influx6/haiku/trees/attrs"
	"github.com/influx6/haiku/trees/elems"
	"github.com/influx6/haiku/trees/events"
	"github.com/influx6/haiku/types"
	"github.com/influx6/haiku/views"
)

func main() {

	page := views.Page()

	var clickMe = func(types.Event) {
		log.Printf("smark down!")
	}

	var menuItem = []string{"shops", "janitor", "booky", "drummer"}

	menu := views.NewView(func(v views.Views) trees.Markup {
		div := elems.Form(elems.Text("welcome"), elems.Paragraph(elems.Label(elems.Text("name")), elems.Input(attrs.Type("text"))))

		var so = elems.Select()
		for _, mi := range menuItem {
			elems.Anchor(
				events.Click(clickMe, "").PreventDefault(),
				attrs.Href("#"+mi),
				elems.Text(mi)).Apply(div)

			so.Augment(elems.Option(attrs.Name(mi), elems.Text(mi)))
		}

		div.Augment(elems.Paragraph(so))
		return div
	})

	page.Mount("body", ".", menu)

	go func() {
		<-time.After(1 * time.Second)
		menuItem = menuItem[1:]
		menu.Send(true)
		<-time.After(1 * time.Second)
		menuItem = append(menuItem, "border", "section", "chief")
		menu.Send(true)
	}()

}
