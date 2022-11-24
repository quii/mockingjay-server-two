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
	http2 "github.com/quii/mockingjay-server-two/domain/mockingjay/http"
)

func main() {
	fs := flag.NewFlagSet("mockingjay", flag.ContinueOnError)

	var (
		adminPort       = fs.String("admin-port", config.DefaultAdminServerPort, "admin server port")
		adminBaseURL    = fs.String("admin-base-url", config.DefaultAdminBaseURL, "admin base url")
		stubPort        = fs.String("stub-port", config.DefaultStubServerPort, "stub server port")
		devMode         = fs.Bool("dev-mode", config.DevModeOff, "dev mode allows templates to be refreshed, logs, etc")
		endpointsFolder = fs.String("endpoints", "", "folder for endpoints")
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

	var endpoints http2.Endpoints
	if endpointsFolder != nil && *endpointsFolder != "" {
		endpoints, err = mockingjay.NewEndpointsFromCue(*endpointsFolder)
		if err != nil {
			log.Fatal(err)
		}
	}

	service, err := mockingjay.NewService(endpoints)
	if err != nil {
		log.Fatal(err)
	}

	stubHandler, adminHandler, err := httpserver.New(service, *adminBaseURL, *devMode)
	if err != nil {
		log.Fatal(err)
	}

	printStartupMessage(endpointsFolder, adminPort, stubPort, adminBaseURL)

	go func() {
		if err := http.ListenAndServe(":"+*adminPort, adminHandler); err != nil {
			log.Fatal(err)
		}
	}()

	if err := http.ListenAndServe(":"+*stubPort, stubHandler); err != nil {
		log.Fatal(err)
	}
}

func printStartupMessage(endpointsFolder *string, adminPort *string, stubPort *string, adminURL *string) {
	executable, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("🚀 mockingjay launched! attempting to listen on %s for admin server, and %s for stub server\n", *adminPort, *stubPort)

	if endpointsFolder != nil && *endpointsFolder != "" {
		log.Printf("📂 endpoints loaded from %s/%s\n", executable, executable+*endpointsFolder)
	}
	log.Printf("💡 visit %s to see the current configuration", *adminURL)
}
