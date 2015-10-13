//Package attr provides attributes for html base elements
//Source Automatically Generated

//go:generate go run generate.go

package attrs

import (
  "github.com/influx6/haiku/trees"
)

// InputType defines the set type of input values for the input elements
type InputType string
  
const (
  
    TypeButton InputType = "button"
    
    TypeCheckbox InputType = "checkbox"
    
    TypeColor InputType = "color"
    
    TypeDate InputType = "date"
    
    TypeDatetime InputType = "datetime"
    
    TypeDatetimelocal InputType = "datetime-local"
    
    TypeEmail InputType = "email"
    
    TypeFile InputType = "file"
    
    TypeHidden InputType = "hidden"
    
    TypeImage InputType = "image"
    
    TypeMonth InputType = "month"
    
    TypeNumber InputType = "number"
    
    TypePassword InputType = "password"
    
    TypeRadio InputType = "radio"
    
    TypeRange InputType = "range"
    
    TypeMin InputType = "min"
    
    TypeMax InputType = "max"
    
    TypeValue InputType = "value"
    
    TypeStep InputType = "step"
    
    TypeReset InputType = "reset"
    
    TypeSearch InputType = "search"
    
    TypeSubmit InputType = "submit"
    
    TypeTel InputType = "tel"
    
    TypeText InputType = "text"
    
    TypeTime InputType = "time"
    
    TypeUrl InputType = "url"
    
    TypeWeek InputType = "week"
    
)
  
// Name defines attributes of type "Name" for html element types
func Name(val string) *trees.Attribute {
  return &trees.Attribute{Name: "name", Value: val }
}
    
// Checked defines attributes of type "Checked" for html element types
func Checked(val string) *trees.Attribute {
  return &trees.Attribute{Name: "checked", Value: val }
}
    
// ClassName defines attributes of type "ClassName" for html element types
func ClassName(val string) *trees.Attribute {
  return &trees.Attribute{Name: "className", Value: val }
}
    
// Autofocus defines attributes of type "Autofocus" for html element types
func Autofocus(val string) *trees.Attribute {
  return &trees.Attribute{Name: "autofocus", Value: val }
}
    
// Id defines attributes of type "Id" for html element types
func Id(val string) *trees.Attribute {
  return &trees.Attribute{Name: "id", Value: val }
}
    
// HtmlFor defines attributes of type "HtmlFor" for html element types
func HtmlFor(val string) *trees.Attribute {
  return &trees.Attribute{Name: "htmlFor", Value: val }
}
    
// Class defines attributes of type "Class" for html element types
func Class(val string) *trees.Attribute {
  return &trees.Attribute{Name: "class", Value: val }
}
    
// Src defines attributes of type "Src" for html element types
func Src(val string) *trees.Attribute {
  return &trees.Attribute{Name: "src", Value: val }
}
    
// Href defines attributes of type "Href" for html element types
func Href(val string) *trees.Attribute {
  return &trees.Attribute{Name: "href", Value: val }
}
    
// Rel defines attributes of type "Rel" for html element types
func Rel(val string) *trees.Attribute {
  return &trees.Attribute{Name: "rel", Value: val }
}
    
// Type defines attributes of type "Type" for html element types
func Type(val string) *trees.Attribute {
  return &trees.Attribute{Name: "type", Value: val }
}
    
// Placeholder defines attributes of type "Placeholder" for html element types
func Placeholder(val string) *trees.Attribute {
  return &trees.Attribute{Name: "placeholder", Value: val }
}
    
// Value defines attributes of type "Value" for html element types
func Value(val string) *trees.Attribute {
  return &trees.Attribute{Name: "value", Value: val }
}
    