package httpserver

import "github.com/quii/mockingjay-server-two/adapters/httpserver/handlers"

type MockingjayServerService interface {
	handlers.StubServerService
	handlers.AdminServiceService
}

func New(service MockingjayServerService, adminBaseURL string, devMode bool) (*handlers.StubHandler, *handlers.AdminHandler, error) {
	adminHandler, err := handlers.NewAdminHandler(service, devMode)
	if err != nil {
		return nil, nil, err
	}
	return handlers.NewStubHandler(service, adminBaseURL), adminHandler, nil
}
