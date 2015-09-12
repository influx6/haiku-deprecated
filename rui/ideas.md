#ideas
 We need to combine vdom and html templating as a means of providing a core means
 of creating reactive ui for go, for example like below with the current example

   ```html
     <div react-slide>{.name}</div>
   ```

 The idea is to be able to use reactive attributes and properties to allow objects
 of specific kinds to be effected for change in behaviour and to allow the data of
 elements to react to change in value.

 Core Reactive Scope:
   - Behavioural/State Reactivity: elements can change behaviour depending on attribute or state

   - Data Reactivity: elements can react to change in data

 Examples:

  - Behavioural Reactivity: Behavioural reactivity should be have effect within the dom alone except for behaviours
  that can be capable of working in both dom and non-dom environment

      an ordinary div:
        ```html
        <div class="cloud"></div>
        ```

      should be able to become behavioural when it recieves an attributal change:

      ```html
        <div class="cloud" vibrate></div>
      ```

  - Data Reactivity: Data reactivity should embodie the change of values and trees of the behaviours of
    elements, where elements should be able to render themselves within a good and usable manner. Where a single
    change in the data values of a element,will cause a change in the rendered data of an element.

     ```go

        RenderInterface{
          Render()
          Dirty()
        }

        Renderers []RenderInterface

        DataTree := Struct{
          Name: Reactive{"alex"}),
          Rates: DataTree({
            Counter: (Reactive{1}),
            Rate: (Reactive{43})
          }),
        }

        //the name-tag for the element to have
        atom = Renderer(domNode,(`
           <div class={{DataTree.name}}>
             <label>Rates: {{DataTree.rates.rate}}</label>
             <label>Count: {{DataTree.rates.counter}}</label>
           </div>
        `)).use(DataTree) /* =>
           <div class=alex>
             <label>Rates: 1</label>
             <label>Count: 43</label>
           </div>
        */


        //Mount sets up the parent to watch for this element type
        atom.Mount(DOMNode(div))

        //bigest issue is it has become normal to do embedd the render code
        //in the template and not as part of the code wonder if we can do This
        //put the template into a larger template file and just pass in the struct
        //but we loose a way to actively update that data directly since we dont
        //generate it from the core target,best to stick with this approach but
        //find a seamless way of integrating the code to render experience to be
        //easy

     ```

  DataTree are structs that can register internal render interfaces who can affect its final rendering
     result to which the result can then be passed down into vdom for its final result to be displayed to the browser or rendering console.

     How do we get the data into the tree:

      - Do we provide a pass to parent and let the children pick of their parts(very useful for dataquery replies)
      - Do we provide cursors and affect the children directly
      - Do we provide a central data repo that gets data and passes it to everyone who picks off its own
      - Do we combine a central data repo and parent who shares data to its cursors
