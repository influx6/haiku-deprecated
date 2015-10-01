package styles

import (
	"strconv"

	"github.com/influx6/haiku/trees"
)

// Size presents a basic stringifed unit
type Size string

// Px returns the value in "%px" format
func Px(pixels int) Size {
	return Size(strconv.Itoa(pixels) + "px")
}

// Color provides the color style value
func Color(value string) *trees.Style {
	return &trees.Style{Name: "color", Value: value}
}

// Height provides the height style value
func Height(size Size) *trees.Style {
	return &trees.Style{Name: "height", Value: string(size)}
}

// Margin provides the margin style value
func Margin(size Size) *trees.Style {
	return &trees.Style{Name: "margin", Value: string(size)}
}

// Width provides the width style value
func Width(size Size) *trees.Style {
	return &trees.Style{Name: "width", Value: string(size)}
}
