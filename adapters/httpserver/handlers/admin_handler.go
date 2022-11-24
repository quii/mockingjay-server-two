package handlers

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/quii/mockingjay-server-two/domain/crud"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/contract"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/stub"
)

var (
	//go:embed "templates/*"
	templates embed.FS

	//go:embed static
	static embed.FS
)

const (
	HeaderMockingjayMatched    = "x-mockingjay-matched"
	HeaderMockingjayMatchID    = "x-mockingjay-match-id"
	ReportsPath                = "/match-reports"
	EndpointsPath              = "/"
	CDCPath                    = "/cdc"
	contentTypeApplicationJSON = "application/json"
)

type AdminServiceService interface {
	Reports() crud.CRUDesque[uuid.UUID, matching.Report]
	Endpoints() crud.CRUDesque[uuid.UUID, stub.Endpoint]
	CheckEndpoints() ([]contract.Report, error)
}

type HTTPRenderer interface {
	Render(w http.ResponseWriter, accept string, template string, thing any)
}

type AdminHandler struct {
	http.Handler
	service AdminServiceService
}

func NewAdminHandler(service AdminServiceService, renderer HTTPRenderer) (*AdminHandler, error) {
	app := &AdminHandler{
		service: service,
	}

	reportHandler := ReportHandler{
		service:  service,
		renderer: renderer,
	}

	endpointHandler := EndpointHandler{
		service:  service,
		renderer: renderer,
	}

	cdcHandler := CDCHandler{
		service:  service,
		renderer: renderer,
	}

	adminRouter := mux.NewRouter()
	adminRouter.HandleFunc(ReportsPath, reportHandler.allReports).Methods(http.MethodGet)
	adminRouter.HandleFunc(ReportsPath, reportHandler.deleteReports).Methods(http.MethodDelete)
	adminRouter.HandleFunc(ReportsPath+"/{reportID}", reportHandler.getReport).Methods(http.MethodGet)

	adminRouter.HandleFunc(EndpointsPath, endpointHandler.allEndpoints).Methods(http.MethodGet)
	adminRouter.HandleFunc(EndpointsPath+"{endpointIndex}", endpointHandler.DeleteEndpoint).Methods(http.MethodDelete)
	adminRouter.HandleFunc(EndpointsPath, endpointHandler.addEndpoint).Methods(http.MethodPost)

	adminRouter.HandleFunc(CDCPath, cdcHandler.checkContracts).Methods(http.MethodGet)

	staticHandler, err := newStaticHandler()
	if err != nil {
		return nil, err
	}
	adminRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticHandler))

	app.Handler = adminRouter
	return app, nil
}

func newStaticHandler() (http.Handler, error) {
	lol, err := fs.Sub(static, "static")
	if err != nil {
		return nil, err
	}
	return http.FileServer(http.FS(lol)), nil
}
