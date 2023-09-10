package route

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/masfuulaji/go-chat/internal/app/handlers"
	"github.com/masfuulaji/go-chat/internal/app/services"
	"github.com/masfuulaji/go-chat/internal/database"
	"github.com/masfuulaji/go-chat/internal/websocket"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
func SetupRoute(mux *mux.Router) {
	pool := websocket.NewPool()
	go pool.Start()

    db, err := database.InitMongoDB()
    if err != nil {
        fmt.Println(err.Error())
    }

    mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
        w.Write([]byte("Hello World!"))
	}).Methods("GET")

    mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		services.ServeWs(pool, w, r)
	}).Methods("GET")

    roomHandler := handlers.NewRoomHandler(services.NewRoomService(db))
    room := mux.PathPrefix("/room").Subrouter()
    room.HandleFunc("", roomHandler.CreateRoom).Methods("POST")
    room.HandleFunc("/{id}", roomHandler.GetRoom).Methods("GET")
    room.HandleFunc("", roomHandler.GetRooms).Methods("GET")
    room.HandleFunc("/{id}", roomHandler.UpdateRoom).Methods("PUT")
    room.HandleFunc("/{id}", roomHandler.DeleteRoom).Methods("DELETE")
}
