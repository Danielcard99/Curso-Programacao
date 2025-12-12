package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Full Cycle Rocks!")
	})

	http.ListenAndServe(":8080", nil)
}

