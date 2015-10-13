// +build js

package main

import (
	"time"

	"github.com/influx6/haiku/trees"
	"github.com/influx6/haiku/trees/attrs"
	"github.com/influx6/haiku/trees/elems"
	"github.com/influx6/haiku/views"
)

func main() {

	page := views.Page()

	var menuItem = []string{"shops", "janitor", "booky", "drummer"}

	menu := views.NewView(func() trees.Markup {
		div := elems.Form(elems.Paragraph(elems.Label(elems.Text("name")), elems.Input(attrs.Type("text"))))

		var so = elems.Select()
		for _, mi := range menuItem {
			so.Augment(elems.Option(attrs.Name(mi), elems.Text(mi)))
		}

		div.Augment(elems.Paragraph(so))
		return div
	})

	page.Mount("body", ".", menu)

	go func() {
		<-time.After(2000 * time.Millisecond)
		menuItem = menuItem[1:]
		menu.Send(true)
		<-time.After(2000 * time.Millisecond)
		menuItem = append(menuItem, "border", "section", "chief")
		menu.Send(true)
	}()

	// window := dom.GetWindow()
	// doc := window.Document()
	// slug := doc.GetElementsByTagName("body")[0].GetElementsByTagName("text")[0]
	//
	// log.Printf("attrs %s and data %s", slug.Attributes(), slug.Dataset())
}
