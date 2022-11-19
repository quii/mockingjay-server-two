package httpserver

import "github.com/quii/mockingjay-server-two/adapters/httpserver/handlers"

type MockingjayServerService interface {
	handlers.StubServerService
	handlers.AdminServiceService
}

func New(service MockingjayServerService, adminBaseURL string, devMode bool) (*handlers.StubHandler, *handlers.AdminHandler) {
	return handlers.NewStubHandler(service, adminBaseURL), handlers.NewAdminHandler(service, devMode)
}
