package views

import (
	"testing"

	"github.com/influx6/flux"
	"github.com/influx6/haiku/trees"
	"github.com/influx6/haiku/trees/attrs"
	"github.com/influx6/haiku/trees/elems"
)

//videoData to be rendered
var videoData = []map[string]string{
	map[string]string{
		"src":  "https://youtube.com/xF5R32YF4",
		"name": "Joyride Lewis!",
	},
	map[string]string{
		"src":  "https://youtube.com/dox32YF4",
		"name": "Wonderlust Bombs!",
	},
}

var treeRenderlen = 246

func TestReactiveView(t *testing.T) {
	videos := NewView(func() trees.Markup {
		dom := elems.Div()
		for _, data := range videoData {
			dom.Augment(elems.Video(
				attrs.Src(data["src"]),
				elems.Text(data["name"]),
			))
		}
		return dom
	})

	bo := videos.RenderHTML()

	if len(bo) != treeRenderlen {
		flux.FatalFailed(t, "Rendered result with invalid length, expected %d but got %d -> \n %s", treeRenderlen, len(bo), bo)
	}

	flux.LogPassed(t, "Rendered result accurated with length %d", treeRenderlen)
}
