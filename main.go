package main

import (
	"net/http"

	"github.com/rinosukmandityo/context-example/server"
)

func main() {
	http.HandleFunc("/search", server.HandleSearch)
	http.ListenAndServe(":8000", nil)
}
