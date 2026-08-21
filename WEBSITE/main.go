package main

import (
	"html/template"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/register", handlerReg)
	http.ListenAndServe(":8080", mux)
}


}
