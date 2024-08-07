package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/masfuulaji/go-chat/internal/database"
	"github.com/masfuulaji/go-chat/internal/route"
	"github.com/spf13/viper"
)

func main() {
	_, err := database.InitMongoDB()
	if err != nil {
		fmt.Println(err.Error())
	}
	mux := mux.NewRouter()
	route.SetupRoute(mux)

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	allowedUrl := viper.GetStringSlice("allowed_url")

	handler := handlers.CORS(
		handlers.AllowedOrigins(allowedUrl),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
		handlers.AllowCredentials(),
	)
	http.Handle("/", handler(mux))
	http.ListenAndServe(":8080", nil)
}
