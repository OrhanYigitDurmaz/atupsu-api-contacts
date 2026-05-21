//go:build windows

package main

import (
	"context"
	"net/http"

	"golang.org/x/sys/windows/svc"
)

const serviceName = "AtupsuAPI"

type atupsuService struct {
	server *http.Server
}

func (s *atupsuService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (bool, uint32) {
	changes <- svc.Status{State: svc.StartPending}

	go s.server.ListenAndServe()

	changes <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}

	for c := range r {
		switch c.Cmd {
		case svc.Interrogate:
			changes <- c.CurrentStatus
		case svc.Stop, svc.Shutdown:
			changes <- svc.Status{State: svc.StopPending}
			s.server.Shutdown(context.Background())
			return false, 0
		}
	}
	return false, 0
}

func runService(server *http.Server) error {
	return svc.Run(serviceName, &atupsuService{server: server})
}

func isInteractive() (bool, error) {
	return svc.IsAnInteractiveSession()
}
