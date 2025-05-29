package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const DefaultShutdownTimeout = 30 * time.Second

type HTTPServer struct {
	srv *http.Server
}

func NewHTTPServer(handler http.Handler, addr string) *HTTPServer {
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return &HTTPServer{
		srv: srv,
	}
}

func (s *HTTPServer) Run() {
	go func() {
		if err := s.srv.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	s.Shutdown()
}

func (s *HTTPServer) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}
