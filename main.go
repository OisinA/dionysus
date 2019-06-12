package main

import (
	"github.com/bwmarrin/lit"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"net/http"

	"dionysus/api"
	"dionysus/services"
)

func main() {
	lit.LogLevel = 3
	lit.Prefix = "dionysus"
	err := services.Setup()
	if err != nil {
		lit.Error(err.Error())
		return
	}
	defer services.Cleanup()
	r := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: true,
		MaxAge: 300,
	})
	r.Use(cors.Handler)
	api.SetSecretKey()
	api := api.API{}
	api.Register(r)
	lit.Info("Starting HTTP server")
	err = http.ListenAndServe(":8070", r)
	if err != nil {
		lit.Error(err.Error())
		return
	}
}
