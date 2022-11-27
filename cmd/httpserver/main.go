package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/peterbourgon/ff/v3"
	"github.com/quii/mockingjay-server-two/adapters/config"
	"github.com/quii/mockingjay-server-two/adapters/httpserver"
	"github.com/quii/mockingjay-server-two/domain/mockingjay"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/contract"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/stub"
)

func main() {
	fs := flag.NewFlagSet("mockingjay", flag.ContinueOnError)

	var (
		adminPort       = fs.String("admin-port", config.DefaultAdminServerPort, "admin server port")
		adminBaseURL    = fs.String("admin-base-url", config.DefaultAdminBaseURL, "admin base url")
		stubPort        = fs.String("stub-port", config.DefaultStubServerPort, "stub server port")
		devMode         = fs.Bool("dev-mode", config.DevModeOff, "dev mode allows templates to be refreshed, logs, etc")
		endpointsFolder = fs.String("endpoints", "", "folder for endpoints")
		cdcMode         = fs.Bool("cdc", config.CDCModeOff, "Run the CDCs")
		_               = fs.String("config", "", "config file (optional)")

		err       error
		endpoints stub.Endpoints

		httpClient = &http.Client{Timeout: 5 * time.Second}
		service    = mockingjay.NewService(contract.NewService(httpClient))
	)

	if err := ff.Parse(fs, os.Args[1:],
		ff.WithEnvVarPrefix("mockingjay"),
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.PlainParser),
	); err != nil {
		log.Fatal(err)
	}

	if endpointsFolder != nil && *endpointsFolder != "" {
		endpoints, err = stub.NewEndpointsFromCue(*endpointsFolder)
		if err != nil {
			log.Fatal(err)
		}
		if err := endpoints.Compile(); err != nil {
			log.Fatal(err)
		}
		for _, endpoint := range endpoints {
			if err := service.Endpoints().Create(endpoint.ID, endpoint); err != nil {
				log.Fatal(err)
			}
		}
	}

	if cdcMode != nil && *cdcMode {
		reports, err := service.CheckEndpoints()
		if err != nil {
			log.Fatal(err)
		}
		_ = json.NewEncoder(os.Stdout).Encode(reports)
		os.Exit(0)
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
