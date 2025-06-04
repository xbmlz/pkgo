package main

import (
	"fmt"

	"github.com/xbmlz/pkgo/conf"
	"github.com/xbmlz/pkgo/ginx"
	"github.com/xbmlz/pkgo/log"
	"github.com/xbmlz/pkgo/server"
)

type Config struct {
	Log struct {
		Level string `json:"level" yaml:"level"`
		File  string `json:"file" yaml:"file"`
	} `json:"log" yaml:"log"`
	Server struct {
		Host string `json:"host" yaml:"host"`
		Port int    `json:"port" yaml:"port"`
	} `json:"server" yaml:"server"`
}

func main() {
	// Load config
	cfg := &Config{}
	err := conf.MustLoad("config.yaml").Parse(cfg)
	if err != nil {
		panic(err)
	}

	// Init logger
	log.InitLogger(
		log.WithLevel(cfg.Log.Level),
		log.WithFile(cfg.Log.File),
	)

	// create gin router
	r := ginx.New()

	// start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Infof("Server start at %s", addr)
	srv := server.NewHTTPServer(addr, r)
	srv.Run()
}
