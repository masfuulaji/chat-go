package route

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/masfuulaji/go-chat/internal/app/handlers"
	"github.com/masfuulaji/go-chat/internal/app/services"
	"github.com/masfuulaji/go-chat/internal/websocket"
)

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
    room.Use(AuthMiddleware)
    room.HandleFunc("", roomHandler.CreateRoom).Methods("POST")
    room.HandleFunc("/{id}", roomHandler.GetRoom).Methods("GET")
    room.HandleFunc("", roomHandler.GetRooms).Methods("GET")
    room.HandleFunc("/{id}", roomHandler.UpdateRoom).Methods("PUT")
    room.HandleFunc("/{id}", roomHandler.DeleteRoom).Methods("DELETE")

    messageHandler := handlers.NewMessageHandler(services.NewMessageService())
    message := mux.PathPrefix("/message").Subrouter()
    message.Use(AuthMiddleware)
    message.HandleFunc("", messageHandler.CreateMessage).Methods("POST")
    message.HandleFunc("/{id}", messageHandler.GetMessage).Methods("GET")
    message.HandleFunc("", messageHandler.GetMessages).Methods("GET")
    message.HandleFunc("/{id}", messageHandler.UpdateMessage).Methods("PUT")
    message.HandleFunc("/{id}", messageHandler.DeleteMessage).Methods("DELETE")
    
    authHandler := handlers.NewAuthHandler(services.NewAuthService())
    auth := mux.PathPrefix("/auth").Subrouter()
    auth.HandleFunc("", authHandler.Login).Methods("POST")
    auth.HandleFunc("/login", authHandler.Login).Methods("POST")
    auth.HandleFunc("/register", authHandler.Register).Methods("POST")
    auth.HandleFunc("/check", authHandler.CheckAuth).Methods("GET")

    userHandler := handlers.NewUserHandler(services.NewUserService())
    user := mux.PathPrefix("/user").Subrouter()
    user.Use(AuthMiddleware)
    user.HandleFunc("", userHandler.CreateUser).Methods("POST")
    user.HandleFunc("/{id}", userHandler.GetUser).Methods("GET")
    user.HandleFunc("", userHandler.GetUsers).Methods("GET")
    user.HandleFunc("/{id}", userHandler.UpdateUser).Methods("PUT")
    user.HandleFunc("/{id}", userHandler.DeleteUser).Methods("DELETE")
}
