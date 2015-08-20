# Flux
A library of interesting structs for the adventurous

##Goodies
  - Stack

    Provide a very simple approach to function binding and stacking, rather than depend on an array of callbacks or other approaches to provide a pubsub and eventful system. Stacks simple combine functions like linkedlists using nodes i.e each function binds to the next one in a seamless fashion and allows you to emit/apply values at any level of the stack to either propagate upwards or downwards or have only that function effected, which provides a nice base for higher class systems eg function reactivity


#Examples

  - Stacks

    Provide 7 basic emission functions,with each providing a flexible pattern,also to allow stacks provide that 'Once' like behaviour the stack function format provides the data given and the Stack receiving the data (you could easily close a stack after receiving the data if needed)

    - Identity:

      Identity just like the matrix operation of the product of a value against an identity matrix returns the same value,here it provides somewhat of an sideeffect pattern where you want to have the stack perform its operation on the data supplied but not affect the returned value and the original value being passed down the chain

       ```

      	master := NewStack(func(data interface{}, _ Stacks) interface{} {
          //data = 20
      		log.Println(data.(int) * 20)
          return nil
      	})

      	slave := master.Stack(func(data interface{}, _ Stacks) interface{} {
          //data = 20
      		log.Println(data.(int) / 20)
          return nil //will be replaced by 20
      	}, true)

        mval := master.Identity(20) //mval = 20 and slave receives no data hits

       ```

    - Isolate:

      Isolate is just what it says, stacks walk alot in chains and there are times you only one to use a specific stack for its effect on a value but do not desire that value to propagate to the rest of the chain(upward or downwards),this provides just that use case

       ```

      	master := NewStack(func(data interface{}, _ Stacks) interface{} {
      		return data.(int) * 20
      	})

      	slave := master.Stack(func(data interface{}, _ Stacks) interface{} {
      		return data.(int) / 20
      	}, true)

        mval := master.Isolate(20) //mval = 400 and slave receives no data hits

       ```

    - Call:

      Call is a mix of Isolation with Identity, depending on where the stack is called,the target stack will apply its effect on the data supplied and return its value but propagate the original data to the next stack. Now why this is useful is when a stack diverts at any point is used to create new stacks. Remember the returned value of the stack being used is the modified value by this and only this stack.

       ```

      	master := NewStack(func(data interface{}, _ Stacks) interface{} {
          //data = 20
          return datal.(int) * 20
      	})

      	slave := master.Stack(func(data interface{}, _ Stacks) interface{} {
          //data = 400
          return data.(int) / 10
      	}, true)

      	slave2 := slave.Stack(func(data interface{}, _ Stacks) interface{} {
          //data = 40
          return data + 50
      	}, true)

        master.Call(20) //returns  400 but slave gets 20 and slave2 gets 1

       ```

    - Apply:

      Apply takes Call method a little further by ensuring that when a value is supplied,it is passed down all connected child stacks till the last, hence the returned value received is actually the value returned by the very last stack in the chain.

       ```

      	master := NewStack(func(data interface{}, _ Stacks) interface{} {
          //data = 20
          return datal.(int) * 20
      	})

      	slave := master.Stack(func(data interface{}, _ Stacks) interface{} {
          //data = 400
          return data.(int) / 10
      	}, true)

      	slave2 := slave.Stack(func(data interface{}, _ Stacks) interface{} {
          //data = 40
          return data + 50
      	}, true)

        master.Apply(20) //returns  90

       ```

    - Lift:

      Lift provides an interesting reverse direction but of same effect as apply. When you need to apply a value from the root of the chains without having to get the root stack,lift takes the given value and passes until it reaches the root(i.e the stack that has no root Stack) which then calls Apply,to fire off the ripple effects

       ```

      	master := NewStack(func(data interface{}, _ Stacks) interface{} {
          //data = 20
          return datal.(int) * 20
      	})

      	slave := master.Stack(func(data interface{}, _ Stacks) interface{} {
          //data = 400
          return data.(int) / 10
      	}, true)

      	slave2 := slave.Stack(func(data interface{}, _ Stacks) interface{} {
          //data = 40
          return data + 50
      	}, true)

        slave.Lift(20) //returns  90

       ```

    - Levitate:

      Levitate provides the reverse ripple effect of the Lift function, Levitate allows a bottom-up mutation effect instead of the standard top-down effect,due to the fact that everything is a linked chain,it takes the returned value of each chains and pass it as the value of the upper chain,in this case until the root chain gets the value. The returned value of the stack used to fire this is the returned value of that particular stack after modification by that same stack

       ```

      	master := NewStack(func(data interface{}, _ Stacks) interface{} {
          //data = 20
          return datal.(int) * 20
      	})

      	slave := master.Stack(func(data interface{}, _ Stacks) interface{} {
          //data = 400
          return data.(int) / 10
      	}, true)

      	slave2 := slave.Stack(func(data interface{}, _ Stacks) interface{} {
          //data = 40
          return data + 50
      	}, true)

        slave2.Levitate(20) //returns  70

       ```

    - LiftApply:

      LiftApply provide a mix on Levitate operations by follow the same principles but returning the root value after it receives the mutated returned value from its previous mutation

       ```

      	master := NewStack(func(data interface{}, _ Stacks) interface{} {
          //data = 7
          return datal.(int) * 20
      	})

      	slave := master.Stack(func(data interface{}, _ Stacks) interface{} {
          //data = 70
          return data.(int) / 10
      	}, true)

      	slave2 := slave.Stack(func(data interface{}, _ Stacks) interface{} {
          //data = 20
          return data + 50
      	}, true)

        slave2.LiftApply(20) //returns  140

       ```
