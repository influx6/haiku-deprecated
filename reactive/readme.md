#Reactive
Reactive is based on the idea of FRP[1] and embodies two distinct but useful ideas to create frp style operations in go. Mutations and Function composition with Function reactive designs allow for a vast application and combining these three principles to create reusable,reactive and functional systems


##Install

    ` go get github.com/influx6/reactive `

##Examples

  - Mutations

    Mutations occur in every system and creating a simple and elegant means of meeting this constraints is key to any functional reactive system.  Using the ideals of simplification. Reactive provides mutations on basic types supported by the go language and simple builds larger constructs of structs,maps or arrays based on these types. This simplifies and allows change at a basic, approchable level.

    e.g

     ```go
      age := Transform(1)

      //changes can be listen to with a callback attached

      age.Get() => // 1

      //change the immutable to get a new one

      age.Set(20)

      age.Get() => //20

     ```
