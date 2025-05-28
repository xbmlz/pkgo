package server

import (
	"context"
	"errors"
	"net/http"
	"time"
)

const DefaultShutdownTimeout = 30 * time.Second

type HTTPServer struct {
	srv *http.Server
}

var _ Server = (*HTTPServer)(nil)

func NewHTTPServer(handler http.Handler, addr string) *HTTPServer {
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return &HTTPServer{
		srv: srv,
	}
}

func (s *HTTPServer) Run() error {
	err := s.srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *HTTPServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultShutdownTimeout)
	defer cancel()
	return s.srv.Shutdown(ctx)
}
