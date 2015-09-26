```js

var ok bool
		document, ok = dom.GetWindow().Document().(dom.HTMLDocument)
		if !ok {
			panic("Could not convert document to dom.HTMLDocument")
		}
		browserSupportsPushState = (js.Global.Get("onpopstate") != js.Undefined) &&
			(js.Global.Get("history") != js.Undefined) &&
			(js.Global.Get("history").Get("pushState") != js.Undefined)

      // watchHash listens to the onhashchange event and calls r.pathChanged when
      // it changes
      func (r *Router) watchHash() {
      	js.Global.Set("onhashchange", func() {
      		go func() {
      			path := getPathFromHash(getHash())
      			r.pathChanged(path, false)
      		}()
      	})
      }

      // watchHistory listens to the onpopstate event and calls r.pathChanged when
      // it changes
      func (r *Router) watchHistory() {
      	js.Global.Set("onpopstate", func() {
      		go func() {
      			r.pathChanged(getPath(), false)
      			if r.ShouldInterceptLinks {
      				r.InterceptLinks()
      			}
      		}()
      	})
      }

      // getPathFromHash returns everything after the "#" character in hash.
      func getPathFromHash(hash string) string {
      	return strings.SplitN(hash, "#", 2)[1]
      }

      // getHash is an alias for js.Global.Get("location").Get("hash").String()
      func getHash() string {
      	return js.Global.Get("location").Get("hash").String()
      }

      // setHash is an alias for js.Global.Get("location").Set("hash", hash)
      func setHash(hash string) {
      	js.Global.Get("location").Set("hash", hash)
      }

      // getPath is an alias for js.Global.Get("location").Get("pathname").String()
      func getPath() string {
      	return js.Global.Get("location").Get("pathname").String()
      }

      // pushState is an alias for js.Global.Get("history").Call("pushState", nil, "", path)
      func pushState(path string) {
      	js.Global.Get("history").Call("pushState", nil, "", path)
      }


      func (r *Router) InterceptLinks() {
	for _, link := range document.Links() {
		href := link.GetAttribute("href")
		switch {
		case href == "":
			return
		case strings.HasPrefix(href, "http://"), strings.HasPrefix(href, "https://"), strings.HasPrefix(href, "//"):
			// These are external links and should behave normally.
			return
		case strings.HasPrefix(href, "#"):
			// These are anchor links and should behave normally.
			// Recall that even when we are using the hash trick, href
			// attributes should be relative paths without the "#" and
			// router will handle them appropriately.
			return
		case strings.HasPrefix(href, "/"):
			// These are relative links. The kind that we want to intercept.
			if r.listener != nil {
				// Remove the old listener (if any)
				link.RemoveEventListener("click", true, r.listener)
			}
			r.listener = link.AddEventListener("click", true, r.interceptLink)
		}
	}
}

```
