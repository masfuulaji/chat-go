package route

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/masfuulaji/go-chat/internal/app/handlers"
	"github.com/masfuulaji/go-chat/internal/app/services"
	"github.com/masfuulaji/go-chat/internal/websocket"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
}
func SetupRoute(mux *mux.Router) {
	pool := websocket.NewPool()
	go pool.Start()


    mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
        w.Write([]byte("Hello World!"))
	}).Methods("GET")

    mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		services.ServeWs(pool, w, r)
	}).Methods("GET")

    roomHandler := handlers.NewRoomHandler(services.NewRoomService())
    room := mux.PathPrefix("/room").Subrouter()
    room.HandleFunc("", roomHandler.CreateRoom).Methods("POST")
    room.HandleFunc("/{id}", roomHandler.GetRoom).Methods("GET")
    room.HandleFunc("", roomHandler.GetRooms).Methods("GET")
    room.HandleFunc("/{id}", roomHandler.UpdateRoom).Methods("PUT")
    room.HandleFunc("/{id}", roomHandler.DeleteRoom).Methods("DELETE")

    messageHandler := handlers.NewMessageHandler(services.NewMessageService())
    message := mux.PathPrefix("/message").Subrouter()
    message.HandleFunc("", messageHandler.CreateMessage).Methods("POST")
    message.HandleFunc("/{id}", messageHandler.GetMessages).Methods("GET")
    message.HandleFunc("", messageHandler.GetMessage).Methods("GET")
    message.HandleFunc("/{id}", messageHandler.UpdateMessage).Methods("PUT")
    message.HandleFunc("/{id}", messageHandler.DeleteMessage).Methods("DELETE")
    
    authHandler := handlers.NewAuthHandler(services.NewAuthService())
    auth := mux.PathPrefix("/auth").Subrouter()
    auth.HandleFunc("/login", authHandler.Login).Methods("POST")
    auth.HandleFunc("/register", authHandler.Register).Methods("POST")
}
