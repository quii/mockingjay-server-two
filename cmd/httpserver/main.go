package main

import (
	"flag"
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
		endpointsFolder = fs.String("endpoints", config.DefaultEndpointsLocation, "folder for endpoints")
		_               = fs.String("config", "", "config file (optional)")
	)

	err := ff.Parse(fs, os.Args[1:],
		ff.WithEnvVarPrefix("mockingjay"),
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.PlainParser),
	)

	if err != nil {
		log.Fatal(err)
	}

	endpoints, err := mockingjay.NewEndpointsFromCue(*endpointsFolder)
	if err != nil {
		log.Fatal(err)
	}

	app := httpserver.New(endpoints)

	printStartupMessage(endpointsFolder, adminPort, stubPort)

	go func() {
		if err := http.ListenAndServe(":"+*adminPort, app.AdminRouter); err != nil {
			log.Fatal(err)
		}
	}()

	if err := http.ListenAndServe(":"+*stubPort, http.HandlerFunc(app.StubHandler)); err != nil {
		log.Fatal(err)
	}
}

func printStartupMessage(endpointsFolder *string, adminPort *string, stubPort *string) {
	executable, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	fullPathOfEndpointsFile := executable + *endpointsFolder

	log.Printf("🚀 mockingjay launched! attempting to listen on %s for admin server, and %s for stub server\n", *adminPort, *stubPort)

	if *endpointsFolder == config.DefaultEndpointsLocation {
		log.Println("‼️  no endpoints specified, loading default examples")
	} else {
		log.Printf("📂 endpoints loaded from %s/%s", executable, fullPathOfEndpointsFile)
	}
}
