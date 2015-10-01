// +build ignore

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/influx6/flux"
)

// input contains sets of input types for input elements
var inputs = []string{
	"button",
	"checkbox",
	"color",
	"date",
	"datetime",
	"datetime-local",
	"email",
	"file",
	"hidden",
	"image",
	"month",
	"number",
	"password",
	"radio",
	"range",
	"min",
	"max",
	"value",
	"step",
	"reset",
	"search",
	"submit",
	"tel",
	"text",
	"time",
	"url",
	"week",
}

//attrs contains set of definable attributes
var attrs = []string{
	"checked",
	"className",
	"autofocus",
	"id",
	"htmlFor",
	"class",
	"src",
	"href",
	"rel",
	"type",
	"placeholder",
	"value",
}

func main() {

	file, err := os.Create("./attrs.gen.go")

	if err != nil {
		panic(err)
	}

	defer file.Close()

	fmt.Fprintf(file, `//Package attr provides attributes for html base elements
//Source Automatically Generated

//go:generate go run generate.go

package attrs

import (
  "github.com/influx6/haiku/trees"
)

// InputType defines the set type of input values for the input elements
type InputType string
  `)

	//write out the const input types
	fmt.Fprintf(file, `
const (
  `)

	for _, inp := range inputs {
		capin := flux.Capitalize(strings.Replace(inp, "-", "", -1))
		fmt.Fprintf(file, `
    Type%s InputType = "%s"
    `, capin, inp)
	}

	fmt.Fprintf(file, `
)
  `)

	//write out the attribubtes
	for _, attr := range attrs {
		capattr := flux.Capitalize(attr)
		fmt.Fprintf(file, `
// %s defines attributes of type "%s" for html element types
func %s(val string) *trees.Attribute {
  return &trees.Attribute{Name: "%s", Value: val }
}
    `, capattr, capattr, capattr, attr)
	}

}
