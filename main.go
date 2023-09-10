package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/masfuulaji/go-chat/internal/route"
	"github.com/rs/cors"
)


func main() {
	mux := mux.NewRouter()
	route.SetupRoute(mux)
    handler := cors.Default().Handler(mux)
	http.ListenAndServe(":8080", handler)
}
