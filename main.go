package main

import (
	_ "embed"
	"log"

	"github.com/cue-exp/cueconfig"
)

type Response struct {
	Status int
	Body   string
}

type Request struct {
	Method string
	Path   string
}

type Endpoint struct {
	Description string
	Request     Request
	Response    Response
}

type Endpoints struct {
	Endpoints []Endpoint
}

var (
	//go:embed schema.cue
	schema []byte
)

func main() {
	var endpoints Endpoints

	if err := cueconfig.Load("endpoints.cue", schema, nil, nil, &endpoints); err != nil {
		log.Fatal(err)
	}

	log.Println(endpoints)
}
