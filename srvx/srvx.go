package srvx

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/xbmlz/pkgo/logx"
)

// Server is transport server.
type Server interface {
	Run() error
	Shutdown() error
}

type server struct {
	servers []Server
}

func New(servers ...Server) *server {
	return &server{servers: servers}
}

func (srv *server) Run() error {
	if len(srv.servers) == 0 {
		return nil
	}

	// Create a context that is canceled on receiving termination signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	errCh := make(chan error, 1)

	for _, s := range srv.servers {
		go func(srv Server) {
			if err := srv.Run(); err != nil {
				logx.Errorf("failed to start server, err: %s", err)
				errCh <- err
			}
		}(s)
	}

	select {
	case err := <-errCh:
		_ = srv.Shutdown()
		return err
	case <-ctx.Done():
		return srv.Shutdown()
	case <-quit:
		return srv.Shutdown()
	}
}

// Stop application stop
func (srv *server) Shutdown() error {
	wg := sync.WaitGroup{}
	for _, s := range srv.servers {
		wg.Add(1)
		go func(srv Server) {
			defer wg.Done()
			if err := srv.Shutdown(); err != nil {
				logx.Errorf("failed to stop server, err: %s", err)
			}
		}(s)
	}
	// wait all server graceful shutdown
	wg.Wait()
	return nil
}
