package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"github.com/masfuulaji/go-chat/internal/database"
	"github.com/masfuulaji/go-chat/internal/route"
)


func main() {
    _, err := database.InitMongoDB()
    if err != nil {
        fmt.Println(err.Error())
    }
	mux := mux.NewRouter()
	route.SetupRoute(mux)
    // handler := cors.Default().Handler(mux)
    handler := handlers.CORS(
        handlers.AllowedOrigins([]string{"http://localhost:5173"}),
        handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
        handlers.AllowedHeaders([]string{"Content-Type"}),
        handlers.AllowCredentials(),
    )
    http.Handle("/", handler(mux))
	http.ListenAndServe(":8080", nil)
}
