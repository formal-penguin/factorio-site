package main

import (
	"fmt"
	"net/http"
)

// put ALL the handlers here
func homeHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "<h1>%s</h1><div>%s</div>", "Hello", "This is a test, with the new improved router")
}

// start the server
func main() {
	port := "8000"
	r := NewRouter()
	r.Add("GET", "/", homeHandler)
	fmt.Printf("Listening on port %s...\nPress more Ctrl-C to exit", port)
	http.ListenAndServe(":"+port, r)
}
