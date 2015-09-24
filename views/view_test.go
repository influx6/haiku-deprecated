package views

import (
	"log"
	"testing"

	"github.com/influx6/flux"
)

func TestView(t *testing.T) {
	//videoData to be rendered
	videoData := []map[string]interface{}{
		map[string]interface{}{
			"src":  "https://youtube.com/xF5R32YF4",
			"name": "Joyride Lewis!",
		},
		map[string]interface{}{
			"src":  "https://youtube.com/dox32YF4",
			"name": "Wonderlust Bombs!",
		},
	}

	videos, err := NewTemplateRenderable(`
    <ul>
      {{ range . }}
        <li>
          <video src="{{.src}}">{{.name}}</video>
        <li>
      {{end}}
    </ul>
  `)

	if err != nil {
		flux.FatalFailed(t, "Unable to create video renderer: %s", err)
	}

	home, err := SourceView("homeView", `
    <html>
      <head></head>
      <body>
        <div class="videos">
          {{ (.View "video").RenderHTML }}
        </video>

        <div class="filesystem">
          {{ (.View "home").RenderHTML }}
        </div>
      </body>
    </html>
  `)

	if err != nil {
		flux.FatalFailed(t, "Unable to create sourceview: %s", err)
	}

	err = videos.Execute(videoData)

	if err != nil {
		flux.FatalFailed(t, "Unable to process video template: %s", err)
	}

	home.AddView("video", "video", videos)

	bo := home.Render()
	log.Printf("render: \n%s", bo)
}
