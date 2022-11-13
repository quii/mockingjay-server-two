package httpserver

type MockingjayServerService interface {
	StubServerService
	AdminServiceService
}

func NewServer(service MockingjayServerService, adminBaseURL string) (*StubHandler, *AdminHandler) {
	return NewStubHandler(adminBaseURL, service), NewAdminHandler(service)
}
