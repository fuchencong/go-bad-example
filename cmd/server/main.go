package main

import (
	"log"
	"os"

	"github.com/fuchencong/go-bad-example/internal/app"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	store := app.NewStore("data/board.json")
	server := app.NewServer(store)
	log.Printf("bad example API listening on http://localhost:%s", port)
	log.Fatal(server.Router().Run(":" + port))
}
