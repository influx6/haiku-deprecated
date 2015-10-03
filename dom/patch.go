package dom

import (
	"fmt"
	"log"
	"strings"

	hodom "honnef.co/go/js/dom"
)

// CreateFragment returns a DocumentFragment with the given html dom
func CreateFragment(html string) hodom.DocumentFragment {
	//if we are not in a browser,panic
	panicBrowserDetect()

	//grab the document
	doc := hodom.GetWindow().Document()

	//we need to use innerhtml but DocumentFragments dont have that so we use
	//a discardable div
	div := doc.CreateElement("div")
	//build up the html right in the div
	div.SetInnerHTML(html)

	//create the document fragment
	fragment := doc.CreateDocumentFragment()

	//add the nodes from the div into the fragment
	for _, node := range div.ChildNodes() {
		fragment.AppendChild(node)
	}

	return fragment
}

// CopyNodes copy nodes from one element into another
func CopyNodes(dest hodom.Element, src []hodom.Node) {
	// add the nodes from the list into the element
	for _, node := range src {
		dest.AppendChild(node)
	}
}

// AddNodeIfNone uses the dom.Node.IsEqualNode method to check if not already exist and if so swap them else just add
// NOTE: bad way of doing it use it as last option
func AddNodeIfNone(dest hodom.Node, src hodom.Node) {
	AddNodeIfNoneInList(dest, dest.ChildNodes(), src)
}

// AddNodeIfNoneInList checks a node in a node list if it finds an equal it replaces only else does nothing
func AddNodeIfNoneInList(dest hodom.Node, against []hodom.Node, with hodom.Node) bool {
	for _, no := range against {
		if no.IsEqualNode(with) {
			dest.ReplaceChild(with, no)
			return true
		}
	}
	//not matching, add it
	dest.AppendChild(with)
	return false
}

// AddElementIfNoneInList checks a node in a node list if it finds an equal it replaces only else does nothing
func AddElementIfNoneInList(dest hodom.Node, against []hodom.Element, with hodom.Node) bool {
	for _, no := range against {
		if no.IsEqualNode(with) {
			dest.ReplaceChild(with, no)
			return true
		}
	}
	//not matching, add it
	dest.AppendChild(with)
	return false
}

// Patch takes a dom string and creates a documentfragment from it and patches a existing dom element that is supplied. This algorithim only ever goes one-level deep, its not performant
func Patch(fragment hodom.Node, live hodom.Element) {
	PatchTree(fragment, live, 0)
}

// PatchTree takes a node or document fragment with a live element.This algorithim only ever goes one-level deep, its not performant to try and check every element. Because we know the fragment will probably only contain the root element that contains the children we will instead factor that in and go a level deep to check if the children recieved any changes by new nodes or hash change and act appropriately
// WARNING: this method is specifically geared to dealing with the haiku.Tree dom generation
func PatchTree(fragment hodom.Node, live hodom.Element, level int) {
	if !live.HasChildNodes() {
		// if the live element is actually empty, then just append the fragment which
		// actually appends the nodes within it efficiently
		live.AppendChild(fragment)
		return
	}

	shadowNodes := fragment.ChildNodes()
	// FIXED: instead of going through the children which may be many,
	// liveNodes := fragment.ChildNodes()

patchloop:
	for _, node := range shadowNodes {
		switch node.(type) {
		case hodom.Element:
			elem := node.(hodom.Element)

			//get the tagname
			tagname := elem.TagName()

			// get the basic attrs
			var id, hash, class, uid string

			// do we have 'id' attribute? if so its a awesome chance to simplify
			if elem.HasAttribute("id") {
				id = elem.GetAttribute("id")
			}

			if elem.HasAttribute("class") {
				id = elem.GetAttribute("class")
			}

			// lets check for the hash and uid, incase its a pure template based script
			if elem.HasAttribute("hash") {
				hash = elem.GetAttribute("hash")
			}

			if elem.HasAttribute("uid") {
				uid = elem.GetAttribute("uid")
			}

			// if we have no id,class, uid or hash, we digress to bad approach of using Node.IsEqualNode
			if allEmpty(id, hash, uid) {
				AddNodeIfNone(live, node)
			}

			// eliminate which ones are empty and try to use the non empty to get our target
			if allEmpty(hash, uid) {
				// is the id empty also then we know class is not or vise-versa
				if allEmpty(id) {
					// class is it and we only want those that match narrowing our set
					no := live.QuerySelectorAll(class)

					// if none found we add else we replace
					if len(no) <= 0 {
						live.AppendChild(node)
					} else {
						// check the available sets and replace else just add it
						AddElementIfNoneInList(live, no, node)
					}

				} else {
					// id is it and we only want one
					no := live.QuerySelector(fmt.Sprintf("#%s", id))

					// if none found we add else we replace
					if no == nil {
						live.AppendChild(node)
					} else {
						live.ReplaceChild(node, no)
					}
				}

				continue patchloop
			}

			// lets use our unique id to check for the element if it exists
			sel := fmt.Sprintf(`%s[uid='%s']`, strings.ToLower(tagname), uid)

			log.Printf("using query for queryselect (%s)", sel)
			// we know hash and uid are not empty so we kick ass the easy way
			target := live.QuerySelector(sel)
			log.Printf("using query for queryselect %s -> %+s - %+s", sel, target, live)

			// if we are nil then its a new node add it and return
			if target == nil {
				live.AppendChild(node)
				continue patchloop
			}

			// if the target hash is exactly the same with ours skip it
			if target.GetAttribute("hash") == hash {
				continue patchloop
			}

			//so we got this dude, are we already one level deep ? if so swap else
			// run through the children with Patch
			if level >= 1 {
				live.ReplaceChild(node, target)
			} else {
				PatchTree(node, target, 1)
			}

		default:
			// add it if its not an element
			live.AppendChild(node)
		}
	}
}

// allEmpty checks if all strings supplied are empty
func allEmpty(s ...string) bool {
	var state = true

	for _, so := range s {
		if so != "" {
			state = false
			return state
		}
	}

	return state
}
