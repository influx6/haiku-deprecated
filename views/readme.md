#Views
 Views provide a simple state-machine based heirarchical rendering of data or templates where depending on the state of activator(generally can be anything that triggers a change state for the machine e.g url, url hash change,..etc) will call up states that match that specific criteria.

## Engine
 Views use a url like state map or state address to define a current state or a state to be transitioned to. This allows views to embed other views which register that own laws into the parent view and leave parent view to delegate and provide the state criteria which allows each state to decide to render or not

### Engine Rendering Rules
A basic demonstration of the view engine and state machine operation. Where the (n: m) shows the relation of view tag and view state address i.e (books: .books) defines a view of tag 'books' with a state address of '.books' the 'root->books' order or in url speak, a '/books' path

      *Examples will be laced and demostrated using html symantics and tag rules we can demostrate as below*


    State trees define the state of the state machine engine and the hierarchical nature of views. Views are locked in by the tag they are giving so a FilesView with tag drop becomes a 'drop' state address point,making it easy to provide same view handlers except when specific to use the root address point

       ```go
        RootView(root: .)
        |-VideoView(video: .)
        |-HomeView(home: .home)
            |-AddressView(address: .)
            |-FilesView(files: .home.files)
            |-FilesView(drops: .home.drops)
            |-FolderView(folders: .home.folders)
                |-FilesViews(files: .home.folders.files)

       ```


 - **Partial Views (using the `.PartialView(string)` function)**
    Partial view entails the single rendering of a view and its sub views within a large hierarchy of views. That is when the state engine receives a state address of a specific subview, that view only renders its self and other subview that are linked as rootviews in it excluding any subview that is a distinct state address from it. Such examples using the state tree above are:


   - **Engine: State->'.'**
          When the engine is set to a '.' state, all views registerd under the rootView with a '.' state address will be active while others who have specific tags like '.home' will not since they do not match the root state address criteria.

          ```html
            /*Rendering:*/

              <rootview>
                <videoview></videoview>
              </rootview>

          ```

   - **Engine: State->'.home'**
          When the engine is set to a '.home' state,all views matching '.home' and any subview matching under home with a '.' root path will be active but any sub view that are specific e.g '.home.files' or '.home.folders' will not be as they are specific in their state address

          ```html
            /*Rendering:*/
              <homeview>
                <addressview></addressview>
              </homeview>
          ```

   - **Full Views(including hierarchical layers: using the .AllView(string) function call)**
    Full views entails the total rendering of a view and its sub views including their hierarchy layering in each of it own views. Just like the Partial views, full view is rendered at the selected root but unlike a partial view that enforces only states matching the state address, this allows a more freedom that every view along the state address paths get rendered down to the last view within that path thereby resolving the returned rendering with each view showing the heirarchy of their levels and layers. Examples of this is describe below:

        - **Engine: State->'.home'**
          When the engine receives a state address ('.home'), it renders the selected root view and the sequentive state path ensuring to respect nesting and hierarchy, where the '.' renders

          ```html
            /*Rendering:*/

              <rootview>
                <videoview></videoview>
                <homeview>
                  <addressview></addressview>
                </homeview>
              </rootview>

          ```

        - **Engine: State->'.home.files'**
           When the engine receives a state address ('.home.files'), it renders the selected root view and the sequentive state path ensuring to respect nesting and hierarchy, where the '.' renders then 'home' renders within the '.' root response and files renders within 'home' response

           ```html
            /*Rendering:*/

              <rootview>
                <videoview></videoview>
                <homeview>
                  <addressview></addressview>
                  <filesview>

                  </filesview>
                </homeview>
              </rootview>

           ```

##Examples

  - Basic View:
    The basic view system built ontop of the state machine is simple, we can create separate views that connect to a large view and only have the currently active view to be present,allowing the system to behave accordingly

    ```go

    videos := View(`
      <ul>
        {{. range}}
          <li>
            <video src="{{.src}}">{{.name}}</video>
          <li>
        {{end}}
      </ul>
    `)

      home := View(`
        <html>
          <head></head>
          <body>
            <div class="videos">
              {{.View('video')}}
            </video>

            <div class="filesystem">
              {{.View('home').View('folders')}}
            </div>
          </body>
        </html>
      `)


     home.UseView("video",videos)
    ```
