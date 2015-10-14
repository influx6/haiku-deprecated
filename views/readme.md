#Haiku views
 Haiku views embrace the core idea of state machines and use in in their own inner workings when interacting with the page, generally you would use this along with the browser location API, but they provide a nice means of managing, how views work and how Haiku defines what should show or should not in a specific instant in time. Generally, the explanations below, provide a colored look into the inner workings and tries to demonstrate how the system functions underneath. But to be clear, views have the right to decide how they will render within a active or deactivated state, i.e the total work of the state machine is to tell the view whether its in active or an inactive state which then the view either decided to render an empty result or a rich content result.

##Internals

  - Engine
    State engine use a url like state map or state address to define a current state or a state to be transitioned to. This allows state views to embed other state views which register that own laws into the parent state view and leave parent state view to delegate and provide the state criteria which allows each state to decide to render or not

    **Engine Rendering Rules**
        A basic demonstration of the state view engine and state machine operation. Where the (n: m) shows the relation of state view tag and state view state address i.e (books: .books) defines a state view of tag 'books' with a state address of '.books' the 'root->books' order or in url speak, a '/books' path

            *Examples will be laced and demostrated using html symantics and tag rules so that there will be more easy understanding of the innerworkings of the state machine and how it will affect the view system, but please not that this examples are to simplify your understanding and not define how views render but how views get activated and deactivated for rendering as they so see fit to do so*


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
