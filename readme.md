# Haiku [[Haiku-ui](https://haiku-ui.io)]
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/influx6/haiku)
[![Travis](https://travis-ci.org/influx6/haiku.svg?branch=master)](https://travis-ci.org/influx6/haiku)

Haiku (although with a interesting name) actually is a ui framework that tries to use simple but effective approaches in the development and management of UI structures/views. By combining multiple short ideas we can create modular and reusable structures that can work both on the backend and frontend with little modification.

# Install

    go get -u github.com/influx6/haiku/...

# Usage
  See [Haiku-ui](https://github.com/influx6/haiku-ui) for a more in depth examples and approaches in using Haiku.

## Goals
  - Create simple go idiomatic view mangement system.
  - Create easy data reactive mechanism that fit the go way.
  - Work either on the server or the client with minimal change (atleast on the view-code side).
  - Be fast and efficient even when compiled down with [Gopherjs](https://github.com/gopherjs/gopherjs).

## Structures
The structures in this library shamefully only connection with the real essense of haiku is that they follow a principle and idea that all parts are made to be small, short and direct to what they do best, i.e no multi-functional piece that tries to perform all the magic in the world but rather each piece fits a particular part of the puzzle and when combined together then creates the immense capability to be used in different ways to achieve their end. Such structures include:

  - **Reactors:** Built on the idea of [FRP](https://en.wikipedia.org/wiki/FRP) functional reactors which is just a big word for functions that can react to change or signals which provide very simple but highly composable structures that can even build allow one to build time-machine like debugging functionality since all signals are atomic and each represent a change of time hence allowing some interesting side effects. But the core idea is we what structs and composable structures that can both react to change and provide signal on change. By combing ideas of reactive functions and reactive atomic signals reactors provide a nice simple way of designing in a go idiomatic approach of dealing with this style of problem solving.

  - **State Engines:** Haiku comes with a very specific,one-view state machine that focuses mainly on the changing of address, why?, because apart from reactive data in views you also want views that can react to the browser change in location, some implement it directly on the pushState and popState of the browser, others using the hash # change while others combine the two for a hybrid solution, Haiku state machines is a unified answer to this at times fractal solutions. The fact is be it either browser location change either by the back or forward button or by the pop and push state mechanisms, all of this are simply signals inherently similar to the very principles of FRP, thereby by creating a state machine able to consume this district signals of change which we call "state address" in haiku state machine speak, we can create a system that can both react to browser view signals and render according. What if all you had to do was tell you views the address that change and the views not included within those address effectively hide themselves and the real focus groups render? That thinking and that idea fueled the need for such a simple but unified mechanism. But its not just on the frontend alone, this system works on the backend because from start, Haiku was meant to be agnostic as towards the area it works.

## Current Features
  - State based view management using a state machine to react to address changes
  - Nestable views
  - Pure string rendering of views using either go templates or custom markup system (this provide a better option)
  - DOM elements tied to Haiku.Components capable of rendering views (using a documentfragment patching method)
  - Blueprinting the creation of components to easily rebuild them

## Future Features
  - Create a simple virtual dom diffing system that reduces rendering overhead
  - Allow server side rendering of content that can be dehydrated on the client with minimal effort

## Dev
- Go Version: 1.5
- Vendor Management: Glide (https://github.com/Masterminds/glide)
