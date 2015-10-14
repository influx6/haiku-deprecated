# Haiku [[Haiku-ui](https://github.com/influx6/haiku-ui)]
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/influx6/haiku)
[![Travis](https://travis-ci.org/influx6/haiku.svg?branch=master)](https://travis-ci.org/influx6/haiku)

Haiku (although with a interesting name) actually is a component rendering framework that tries to use simple but effective approaches in the in how you render your components or views. Its combines virtual dom diffing and state machine driving views to provide a simple API that allows you render any where whether it be on the server or on the client.

# Install

    go get -u github.com/influx6/haiku/...

# Usage
  See [Haiku-ui](https://github.com/influx6/haiku-ui) for a more in depth examples and approaches in using Haiku.

## Goals
  - Create simple go idiomatic view.
  - Create simple reactive views capable of plugging into each other for notification of change.
  - Work either on the server or the client with minimal change (atleast on the view-code side).
  - Be fast and efficient even when compiled down with [Gopherjs](https://github.com/gopherjs/gopherjs).

## Current Features
  - State based view management using a state machine to react to address changes
  - Pure string rendering of views using custom markup system
  - DOM rendering using custom virtual dom and patching/diffing mechanism
  - Server side rendering of views markup embedded within standard go templates

## Future Features
  - Allow server side synchronization of views using dehydration and hydration of data

## Dev
- Go Version: 1.5
