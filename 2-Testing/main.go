package main

import (
	"2-Testing/myapp"
	"net/http"
)

func main() {
	http.ListenAndServe(":3000", myapp.NewHttpHandler())
}
