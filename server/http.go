package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/xbmlz/pkgo/utils"
)

type Config struct {
	Host        string `yaml:"host" json:"host" mapstructure:"host" env:"HTTP_HOST"`
	Port        int    `yaml:"port" json:"port" mapstructure:"port" env:"HTTP_PORT"`
	ReadTimeout int    `yaml:"read_timeout" json:"read_timeout" mapstructure:"read_timeout" env:"HTTP_READ_TIMEOUT"`
}

func Run(cfg Config, handler http.Handler) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:        addr,
		Handler:     handler,
		ReadTimeout: time.Duration(utils.OrElse(cfg.ReadTimeout, 5)) * time.Second,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}
