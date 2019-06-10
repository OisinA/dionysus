package main

import (
	"github.com/bwmarrin/lit"
	"github.com/go-chi/chi"
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
	api := api.API{}
	api.Register(r)
	lit.Info("Starting HTTP server")
	http.ListenAndServe(":8070", r)
}
