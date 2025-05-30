package httpx

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/xbmlz/pkgo/srvx"
)

const DefaultShutdownTimeout = 30 * time.Second

var _ srvx.Server = (*HTTPServer)(nil)

type HTTPServer struct {
	srv *http.Server
}

func New(addr string, handler http.Handler) *HTTPServer {
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return &HTTPServer{
		srv: srv,
	}
}

func (s *HTTPServer) Run() (err error) {
	err = s.srv.ListenAndServe()
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
