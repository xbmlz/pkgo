package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/xbmlz/pkgo/log"
)

type HTTPServer struct {
	srv *http.Server
}

func NewHTTPServer(addr string, handler http.Handler) *HTTPServer {
	return &HTTPServer{
		srv: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (hs *HTTPServer) Run() {
	go func() {
		if err := hs.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := hs.srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}
