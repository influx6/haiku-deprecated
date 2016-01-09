# Haiku
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/influx6/haiku)

Haiku is a view rendering framework(the V in MVC), built to render on the browser
or server with little code changes as possible. Built on the solid foundation that [Gopherjs](https://github.com/gopherjs/gopherjs) provides. Haiku combines a virtual diffing system to ensure the minimum work done in updating rendered elements and allows the freedom to decide how your data gets into the view.


## Install

    go get -u github.com/influx6/haiku/...


## Example

  ```go

    package main

    import (
        "honnef.co/go/js/dom"
    	"github.com/influx6/haiku/trees"
    	"github.com/influx6/haiku/trees/attrs"
    	"github.com/influx6/haiku/trees/elems"
    )

    type videoList []map[string]string

    func (v videoList) Render(m ...string) trees.Markup {
    	dom := elems.Div()
    	for _, data := range v {
    		dom.Augment(elems.Video(
    			attrs.Src(data["src"]),
    			elems.Text(data["name"]),
    		))
    	}
    	return dom
    }

    func main() {

    	videos := NewView(videoList([]map[string]string{
    		map[string]string{
    			"src":  "https://youtube.com/xF5R32YF4",
    			"name": "Joyride Lewis!",
    		},
    		map[string]string{
    			"src":  "https://youtube.com/dox32YF4",
    			"name": "Wonderlust Bombs!",
    		},
    	}))

    	/*
    videos.RenderHTML() =>
      <div hash="fUPpf3XsV2"  uid="PkOaCB3C" style="">
        <video hash="TzmLsvA7j2"  uid="jq3Xl9gq" src="https://youtube.com/xF5R32YF4" style="">Joyride Lewis!</video>
        <video hash="t8jeXh1JrU"  uid="GSv22Nqb" src="https://youtube.com/dox32YF4" style="">Wonderlust Bombs!</video>
      </div>  
      */

    win := dom.GetWindow()
  	doc := win.Document()

  	body := doc.QuerySelector("body")


  	video.Mount(body.Underlying())


  ```


## Goals
  - Create simple and go idiomatic views.
  - Be fast and efficient even when compiled down with [Gopherjs](https://github.com/gopherjs/gopherjs).

## Current Features
  - State based view management using a state machines.
  - Pure string rendering of views with custom markup structures.
  - DOM Patching/Diffing mechanism
