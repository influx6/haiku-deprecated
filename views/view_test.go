package views

import (
	"testing"

	"github.com/influx6/flux"
	"github.com/influx6/haiku/trees"
	"github.com/influx6/haiku/trees/attrs"
	"github.com/influx6/haiku/trees/elems"
)

var treeRenderlen = 273

type videoList struct {
	lists []map[string]string
}

func (v *videoList) Render(m ...string) trees.Markup {
	dom := elems.Div()
	for _, data := range v.lists {
		dom.Augment(elems.Video(
			attrs.Src(data["src"]),
			elems.Text(data["name"]),
		))
	}
	return dom
}

func TestReactiveView(t *testing.T) {
	videos := NewView(&videoList{[]map[string]string{
		map[string]string{
			"src":  "https://youtube.com/xF5R32YF4",
			"name": "Joyride Lewis!",
		},
		map[string]string{
			"src":  "https://youtube.com/dox32YF4",
			"name": "Wonderlust Bombs!",
		},
	}})

	bo := videos.RenderHTML()

	if len(bo) != treeRenderlen {
		flux.FatalFailed(t, "Rendered result with invalid length, expected %d but got %d -> \n %s", treeRenderlen, len(bo), bo)
	}

	flux.LogPassed(t, "Rendered result accurated with length %d", treeRenderlen)
}
