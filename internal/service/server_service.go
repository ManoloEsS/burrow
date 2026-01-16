package service

type serverService struct {
	currentStatus ServerStatus
	crashDetected bool
	serverCancel  chan struct{}
}

func NewServerService() ServerService {
	return &serverService{
		currentStatus: ServerStatus{
			Running: false,
			Status:  "Server not running",
		},
		serverCancel: make(chan struct{}),
	}
}

func (s *serverService) StartServer(path string) error {
	return nil
}

func (s *serverService) StopServer() error {
	return nil
}

func (s *serverService) HealthCheck() ServerStatus {
	return s.currentStatus
}

func (s *serverService) runServer(path string) {
}

func (s *serverService) handleServerCrash() {
}
