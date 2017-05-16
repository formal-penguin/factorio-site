package main

import (
	"fmt"
	"net/http"
	"strings"
)

var httpVerbs = [...]string{"GET", "POST", "DELETE", "PATCH", "PUT", "HEAD", "OPTIONS", "CONNECT"}

func methodAllowed(input string) bool {
	input = strings.ToUpper(input)
	for _, b := range httpVerbs {
		if input == b {
			return true
		}
	}
	return false
}

type routeHandler func(http.ResponseWriter, *http.Request)

type routeHandle map[string]routeHandler

func newRouteHandle() routeHandle {
	thisMap := make(map[string]routeHandler)
	return thisMap
}

type node struct {
	children  []*node
	component string
	leaf      bool
	handle    routeHandle
}

func (n *node) findChild(component string) (index int, found bool) {
	for i, c := range n.children {
		if component == c.component {
			return i, true
		}
	}
	return -1, false
}

type Router struct {
	root *node
}

func NewRouter() *Router {
	routerRoot := node{children: nil, component: "", leaf: true, handle: make(routeHandle)}
	// routerRoot := node{children: make([]*node, 16), component: "", leaf: true, handle: make(routeHandle)}
	routerRoot.handle["GET"] = defaultRootHandler
	r := Router{root: &routerRoot}
	return &r
}

// put ALL the handlers here
func defaultRootHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "<h1>%s</h1><div>%s</div>", "Hello", "This is a test, with the new improved router")
}

func (r *Router) Add(method, path string, handleFunc routeHandler) {
	if methodAllowed(method) {
		currentNode := r.root
		for _, pathComponent := range strings.Split(path, "/") {
			fmt.Printf("***************** %s **************\nPress Ctrl-C to exit", currentNode.component)
			if pathComponent != "" {
				fmt.Printf("&&&&&&&&&&&&& %s &&\nPress Ctrl-C to exit", currentNode.component)
				if childNodeIndex, found := currentNode.findChild(pathComponent); found {
					fmt.Printf("if")
					currentNode = currentNode.children[childNodeIndex]
				} else {
					fmt.Printf("else") // regardless if this node is a leaf or not, we want it to NOT be one moving forward, so we just set its left property to false here
					currentNode.leaf = false
					newChild := node{children: make([]*node, 8), component: "", leaf: true, handle: make(routeHandle)}
					currentNode.children = append(currentNode.children, &newChild)
					currentNode = &newChild
				}
			}
		}
		// at this point we should be done iterating thru the path and we are at the leaf that we need to attach a handler here
		// if handleFunc, ok := (currentNode.handle)[strings.ToUpper(method)]; ok {
		// 	panic(fmt.Sprintf("%s is already defined for this path.\nLog a bug if you think this meesage is shown in mistake and provide as much info as you can. Thank you.", strings.ToUpper(method)))
		// } else {
		// 	fmt.Printf("***************** %s **************\nPress Ctrl-C to exit", method)
		// 	fmt.Printf("***************** %s **************\nPress Ctrl-C to exit", currentNode.component)
		// 	// fmt.Printf("%s...\nPress Ctrl-C to exit", handleFunc)
		// 	currentNode.handle[strings.ToUpper(method)] = handleFunc
		// }
		fmt.Printf("***************** %s **************\nPress Ctrl-C to exit", method)
		fmt.Printf("***************** %s **************\nPress Ctrl-C to exit", currentNode.component)
		// fmt.Printf("%s...\nPress Ctrl-C to exit", handleFunc)
		currentNode.handle[strings.ToUpper(method)] = handleFunc

	} else {
		panic(fmt.Sprintf("%s is not a valid http verb.\nLog a bug if you think this meesage is shown in mistake and provide as much info as you can. Thank you.", method))
	}
}

func (r *Router) Find(method, path string) (routeHandler, string) {
	// 404(not found) and 405(method not implemented) I think, 300's should be implemented in each individual function after the trie is traversed and found something
	found := true
	currentNode := r.root
	for _, pathComponent := range strings.Split(path, "/") {
		if childNodeIndex, found := currentNode.findChild(pathComponent); found {
			currentNode = currentNode.children[childNodeIndex]
		} else {
			found = false
		}
	}
	if !found {
		return nil, "404"
	} else {
		if handleFunc, ok := (currentNode.handle)[strings.ToUpper(method)]; ok {
			return handleFunc, "what"
		} else {
			return nil, "405"
		}
	}
}

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if handleToRun, statusCode := r.Find("GET", req.URL.Path); handleToRun != nil {
		handleToRun(res, req)
	} else {
		fmt.Fprintf(res, "<h1>%s</h1><div>HTTP Code: %s</div>", "Error", statusCode)
	}
}
