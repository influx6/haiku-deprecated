package main

import (
	"fmt"
	"html/template"
	"strconv"

	hdom "github.com/influx6/haiku/dom"
	"github.com/influx6/haiku/reactive"
	"github.com/influx6/haiku/views"
	dom "honnef.co/go/js/dom"
)

type contact struct {
	Name     reactive.Observers
	Location reactive.Observers
	Tel      reactive.Observers
	Age      reactive.Observers
}

//go:generate gopherjs -m  build main.go -o ./app.js

func main() {

	win := dom.GetWindow().(dom.Window)

	doc := win.Document().(dom.HTMLDocument)

	container := doc.QuerySelector(".container")

	formtemplate := template.Must(template.New("form").Parse(doc.QuerySelector(".contact-form-template").InnerHTML()))
	viewtemplate := template.Must(template.New("view").Parse(doc.QuerySelector(".contact-view-template").InnerHTML()))

	dude := &contact{
		Name:     reactive.ObserveAtom("John", false),
		Location: reactive.ObserveAtom("Lagos,Nigeria", false),
		Tel:      reactive.ObserveAtom("+01 4234534", false),
		Age:      reactive.ObserveAtom(20, false),
	}

	contactTempler, err := reactive.ExecTempler(viewtemplate, dude)

	if err != nil {
		panic(fmt.Sprintf("Contact template error: %s", err))
	}

	contactView, err := views.ReactiveSourceView("contacts.view", `
			{{ (.View "view").RenderHTML }}
	`, nil)

	contactView.AddView("view", ".", contactTempler)

	formTempler, err := reactive.ExecTempler(formtemplate, dude)

	if err != nil {
		panic(fmt.Sprintf("Contact template error: %s", err))
	}

	// _ = container
	formView, err := views.SourceView("contacts.view", `
			{{ (.View "form").RenderHTML }}
	`, nil)

	if err != nil {
		panic(fmt.Sprintf("Contact view error: %s", err))
	}

	formView.AddView("form", ".", formTempler)

	contacts := hdom.NewElement(container.QuerySelector(".contactView"))

	contacts.AddEvent("keyup", ".contact-form p .input").Next(func(ev dom.Event, next views.NextHandler) {
		target := ev.Target()
		name := target.GetAttribute("name")

		var input *dom.HTMLInputElement
		var ok bool

		if input, ok = target.(*dom.HTMLInputElement); !ok {
			return
		}

		switch name {
		case "name":
			dude.Name.Set(input.Value)
		case "age":
			if age, err := strconv.Atoi(input.Value); err == nil {
				dude.Age.Set(age)
			}
		case "location":
			dude.Location.Set(input.Value)
		case "telephone":
			dude.Tel.Set(input.Value)
		}

		next(ev)
	})

	hdom.NewViewElement(contacts, formView, ".form")
	hdom.NewViewElement(contacts, contactView, ".view")

	lists := views.NewBlueprint("ListItems", template.Must(template.New("root").Parse(`
		{{ define "root" }}
			{{ (.Binding) }} - {{ template "extra" . }}
		{{ end }}
	`)))

	li, err := lists.AndView(dude.Name, views.HiddenNameStrategy("root"), template.Must(template.New("extra").Parse(`
		{{ define "extra" }}
			{{ .Get "tag" }}
		{{ end }}
	`)))

	li.Set("tag", "Sunday")

	hdom.ViewComponent(contacts, li, ".nameview")

}
