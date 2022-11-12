package main

import (
	"errors"
	"flag"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/peterbourgon/ff/v3"
	"github.com/quii/mockingjay-server-two/adapters/config"
	"github.com/quii/mockingjay-server-two/adapters/httpserver"
	"github.com/quii/mockingjay-server-two/domain/mockingjay"
)

func main() {
	fs := flag.NewFlagSet("mockingjay", flag.ContinueOnError)

	var (
		adminPort       = fs.String("admin-port", config.DefaultAdminServerPort, "admin server port")
		stubPort        = fs.String("stub-port", config.DefaultStubServerPort, "stub server port")
		endpointsFolder = fs.String("endpoints", "examples/", "folder for endpoints")
		_               = fs.String("config", "", "config file (optional)")
	)

	err := ff.Parse(fs, os.Args[1:],
		ff.WithEnvVarPrefix("mockingjay"),
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(CueConfigLoader),
	)

	if err != nil {
		log.Fatal(err)
	}

	endpoints, err := mockingjay.NewEndpointsFromCue(*endpointsFolder)
	if err != nil {
		log.Fatal(err)
	}

	app := httpserver.New(endpoints)

	log.Printf("🚀 mockingjay launched! attempting to listen on %s for admin server, and %s for stub server", *adminPort, *stubPort)
	log.Printf("📂 endpoints loaded from %s", *endpointsFolder)

	go func() {
		if err := http.ListenAndServe(":"+*adminPort, app.AdminRouter); err != nil {
			log.Fatal(err)
		}
	}()

	if err := http.ListenAndServe(":"+*stubPort, http.HandlerFunc(app.StubHandler)); err != nil {
		log.Fatal(err)
	}
}

func CueConfigLoader(r io.Reader, set func(name, value string) error) error {
	return errors.New("not implemented yet! https://github.com/quii/mockingjay-server-two/issues/10")
}
