package views

import (
	"log"
	"testing"

	"github.com/influx6/flux"
)

var viewableSize = 475

func TestViewWithViewable(t *testing.T) {
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

	//Lets add this as a viewable with the address set to render with the rootView (note: videos is not a StatefulViewable so the address wont matter but when dealing with a StatefulViewable it does,because it determines the visibility of the element)
	home.AddView("video", ".", videos)

	home.Engine().All(".", "")

	bo := home.Render()

	log.Printf("bo: %d", len(bo))
	if len(bo) != viewableSize {
		flux.FatalFailed(t, "Rendered result with invalid length, expected %d but got %d", viewableSize, len(bo))
	}

	flux.LogPassed(t, "Rendered result accurated with length %d", viewableSize)
}

func TestViewWithStatefulViewable(t *testing.T) {
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

	err = videos.Execute(videoData)

	if err != nil {
		flux.FatalFailed(t, "Unable to process video template: %s", err)
	}

	vidom, err := SourceView("videoView", `
    <video-view>
      {{ (.View "video").RenderHTML }}
    </video-view>
  `)

	if err != nil {
		flux.FatalFailed(t, "Unable to create vidom source view: %s", err)
	}

	vidom.AddView("video", ".", videos)

	//create another source view
	adom, err := SourceView("audioView", `
    <audio-view>
     <audio-item src="../mike/sosm.mp3">Mike Rogers: Sound of Snow</audio-item>
     <audio-item src="../nick/ph.mp3">Nickebacks: Photographs</audio-item>
    </audio-view>

  	{{range .Views }}
  			{{ .RenderHTML }}
  	{{ end }}
  `)

	if err != nil {
		flux.FatalFailed(t, "Unable to create audio source view: %s", err)
	}

	index, err := SourceView("indexView", `
  <html>
    <head></head>
    <body>
      <section class="gopher-vids">
            {{ (.View "vom").RenderHTML }}
      </section>

      <section class="fav-audios">
            {{ (.View "adom").RenderHTML }}
      </section>
    </body>
  </html>
  `)

	if err != nil {
		flux.FatalFailed(t, "Unable to create index source view: %s", err)
	}

	//we are adding a StatefulViewable (vidom) so the address     we give here is very important as it defines the scope of the view and when that scope will be active or not active
	index.AddView("vom", "videos", vidom)

	//lets add the audio-view as a subroot view so it renders with the root '.'
	index.AddView("adom", ".", adom)

	//lets first render with the state address for '.'
	rootRender := index.Render(".")

	if len(rootRender) < 100 {
		flux.FatalFailed(t, "Bad render received: %s", rootRender)
	}

	//lets first render with the state address for '.'
	videoRender := index.Render(".videos")

	if len(videoRender) < 100 {
		flux.FatalFailed(t, "Bad render received: %s", videoRender)
	}

	flux.LogPassed(t, "Successfully render states")
}
