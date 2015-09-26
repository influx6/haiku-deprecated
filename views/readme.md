#state views
 state views provide a simple state-machine based heirarchical rendering of data or templates where depending on the state of activator(generally can be anything that triggers a change state for the machine e.g url, url hash change,..etc) will call up states that match that specific criteria.

##Internals

  - Engine
    State engine use a url like state map or state address to define a current state or a state to be transitioned to. This allows state views to embed other state views which register that own laws into the parent state view and leave parent state view to delegate and provide the state criteria which allows each state to decide to render or not

    **Engine Rendering Rules**
        A basic demonstration of the state view engine and state machine operation. Where the (n: m) shows the relation of state view tag and state view state address i.e (books: .books) defines a state view of tag 'books' with a state address of '.books' the 'root->books' order or in url speak, a '/books' path

            *Examples will be laced and demostrated using html symantics and tag rules so that there will be more easy understanding of the innerworkings of the state machine and how it will affect the view system*


          State trees define the state of the state machine engine and the hierarchical nature of state views. Using tags (eg 'home') and address points we will demonstrate how the machine works

             ```go

              Rootstate(root: .)
              |-Videostate(video: .)
              |-Homestate(home: .home)
                  |-Addressstate(address: .)
                  |-Filesstate(files: .home.files)
                  |-Filesstate(drops: .home.drops)
                  |-Folderstate(folders: .home.folders)
                      |-Filesstate(files: .home.folders.files)

             ```


       - **Partial States (using the `.Partial(string)` function)**
          Partial state entails the single rendering of a  state and its sub states within a large hierarchy of state views. That is when the state engine receives a state address of a specific substate view eg '.home', the engine state only activates its that state (in this case 'home') and other substates(only states with address point '.') that are linked as root substate wthin it excluding other states that have distinct state address eg. '.home.files'. For more demonstrated examples using the state tree above and html symantics:


         - **Engine: State->'.'**
                When the engine is set to '.' (aka root state), all state views registered under the root state as sub root states (i.e  with a '.' state address) will be activated while others who have specific state address like '.home' will not since they do not match the root state ('.') address criteria.

                Using the state tree above, we can render the behaviour as

                ```html
                  /*Rendering:*/

                    <rootstate view>
                      <videostate view></videostate view>
                    </rootstate view>

                ```

         - **Engine: State->'.home'**
                When the engine is set to a '.home' state,all state views matching '.home' and any substate view matching under home with a '.' root path will be active but any sub state view that are specific e.g '.home.files' or '.home.folders' will not be as they are specific in their state address

                Using the state tree above, we can render the behaviour as:

                ```html
                  /*Rendering:*/
                    <homestate view>
                      <addressstate view></addressstate view>
                    </homestate view>
                ```

         - **Full States(including hierarchical layers: using the .All(string) function call)**
          Full state views entails the total rendering of a state view and its sub state views including their hierarchy layering in each of it own state views. Just like the Partial state views, full state view is rendered at the selected root but unlike a partial state view that enforces only states matching the state address, this allows a more freedom that every state view along the state address paths get rendered down to the last state view within that path thereby resolving the returned rendering with each state view showing the heirarchy of their levels and layers. Examples of this is describe below:

              - **Engine: State->'.home'**
                When the engine receives a state address ('.home'), it renders the selected root state view and the sequentive state path ensuring to respect nesting and hierarchy, where the '.' renders

                Using the state tree above, we can render the behaviour as

                ```html
                  /*Rendering:*/

                    <rootstate view>
                      <videostate view></videostate view>
                      <homestate view>
                        <addressstate view></addressstate view>
                      </homestate view>
                    </rootstate view>

                ```

              - **Engine: State->'.home.files'**
                 When the engine receives a state address ('.home.files'), it renders the selected root state view and the sequentive state path ensuring to respect nesting and hierarchy, where the '.' renders then 'home' renders within the '.' root response and files renders within 'home' response

                Using the state tree above, we can render the behaviour as

                 ```html
                  /*Rendering:*/

                    <rootstate view>
                      <videostate view></videostate view>
                      <homestate view>
                        <addressstate view></addressstate view>
                        <filesstate view>

                        </filesstate view>
                      </homestate view>
                    </rootstate view>

                 ```

##Examples

  - Basic state view:
    The basic state view system built ontop of the state machine is simple, we can create separate state views that connect to a large state view and only have the currently active state view to be present,allowing the system to behave accordingly

    ```go

        //Mixing Views and PseudoViews together

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

      	videos, _ := NewTemplateRenderable(`
          <ul>
            {{ range . }}
              <li>
                <video src="{{.src}}">{{.name}}</video>
              <li>
            {{end}}
          </ul>
        `)

      	home, _ := Sourcestate view("homestate view", `
          <html>
            <head></head>
            <body>
              <div class="videos">
                {{ (.state view "video").RenderHTML }}
              </video>

              <div class="filesystem">
                {{ (.state view "home").RenderHTML }}
              </div>
            </body>
          </html>
        `)

      	videos.Execute(videoData)

      	home.Addstate view("video", "video", videos)

        //Views are stateful and you need to either directly supply the state address or use a mechanism that supplies it eg. browser urls else you get an empty shell of the template
      	home.Render(".")  /* =>
        <html>
          <head></head>
          <body>
            <div class="videos">

        <ul>

            <li>
              <video src="https://youtube.com/xF5R32YF4">Joyride Lewis!</video>
            <li>

            <li>
              <video src="https://youtube.com/dox32YF4">Wonderlust Bombs!</video>
            <li>

        </ul>

            </video>

            <div class="filesystem">
              Render.Error: "home" state view not found!
            </div>
          </body>
        </html>

  */

    ```

    ```go

    //Mixing and combine Views and PseudoViews

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

    	videos, _ := NewTemplateRenderable(`
        <ul>
          {{ range . }}
            <li>
              <video src="{{.src}}">{{.name}}</video>
            <li>
          {{end}}
        </ul>
      `)

    	_ = videos.Execute(videoData)

    	vidom, _ := SourceView("videoView", `
        <video-view>
          {{ (.View "video").RenderHTML }}
        </video-view>
      `)

    	vidom.AddView("video", ".", videos)

    	//create another source view
    	adom, _ := SourceView("audioView", `
        <audio-view>
         <audio-item src="../mike/sosm.mp3">Mike Rogers: Sound of Snow</audio-item>
         <audio-item src="../nick/ph.mp3">Nickebacks: Photographs</audio-item>
        </audio-view>

      	{{range .Views }}
      			{{ .RenderHTML }}
      	{{ end }}
      `)

    	index, _ := SourceView("indexView", `
      <html>
        <head></head>
        <body>
          <section>
            using added view functions
            {{ view "vom" }}
          </section>

          <section class="gopher-vids">
                {{ (.View "vom").RenderHTML }}
          </section>

          <section class="fav-audios">
                {{ (.View "adom").RenderHTML }}
          </section>
        </body>
      </html>
      `)

      //add videoview as a sub-state view
    	index.AddView("vom", "videos", vidom)
      //add audioView as a subroot view
    	index.AddView("adom", ".", adom)

    	//lets first render with the state address for '.'
      //Views are stateful and you need to either directly supply the state address or use a mechanism that supplies it eg. browser urls else you get an empty shell of the template
    	rootRender := index.Render(".")  /* =>
        <html>
          <head></head>
          <body>
            <section class="gopher-vids">

            </section>

            <section class="fav-audios">

              <audio-view>
               <audio-item src="../mike/sosm.mp3">Mike Rogers: Sound of Snow</audio-item>
               <audio-item src="../nick/ph.mp3">Nickebacks: Photographs</audio-item>
              </audio-view>

            </section>
          </body>
        </html>
      */


    	//lets render with the state address for '.videos'
      //Views are stateful and you need to either directly supply the state address or use a mechanism that supplies it eg. browser urls else you get an empty shell of the template
    	videoRender := index.Render(".videos") /* =>
      <html>
        <head></head>
        <body>
          <section class="gopher-vids">

          <video-view>
            <ul>
                <li>
                  <video src="https://youtube.com/xF5R32YF4">Joyride Lewis!</video>
                <li>
                <li>
                  <video src="https://youtube.com/dox32YF4">Wonderlust Bombs!</video>
                <li>
            </ul>
          </video-view>

          </section>

            <section class="fav-audios">

              <audio-view>
               <audio-item src="../mike/sosm.mp3">Mike Rogers: Sound of Snow</audio-item>
               <audio-item src="../nick/ph.mp3">Nickebacks: Photographs</audio-item>
              </audio-view>

            </section>
        </body>
      </html>      
      */


    ```
