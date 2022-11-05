package main

import (
	"log"
	"net/http"

	"github.com/quii/mockingjay-server-two/adapters/httpserver"
	"github.com/quii/mockingjay-server-two/domain/endpoints"
)

func main() {
	var endpoints endpoints.Endpoints

	go func() {
		if err := http.ListenAndServe(":8081", httpserver.ConfigServer(&endpoints)); err != nil {
			log.Fatal(err)
		}
	}()

	if err := http.ListenAndServe(":8080", httpserver.NewStubServer(&endpoints)); err != nil {
		log.Fatal(err)
	}
}
