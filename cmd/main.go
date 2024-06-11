package main

import (
	"log"
	"net/http"

	"github.com/sundayezeilo/pismo/app"
	"github.com/sundayezeilo/pismo/config"
)

func main() {
	cfg := config.LoadEnv("")
	app := app.NewApp(cfg)
	defer app.Cleanup()

	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      app.Router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	log.Printf("Server started on port %v", server.Addr)
	log.Fatal(server.ListenAndServe())
}
