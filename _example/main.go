package main

import (
	"fmt"

	"github.com/xbmlz/pkgo/confx"
	"github.com/xbmlz/pkgo/cronx"
	"github.com/xbmlz/pkgo/ginx"
	"github.com/xbmlz/pkgo/logx"
	"github.com/xbmlz/pkgo/srvx"
	"github.com/xbmlz/pkgo/srvx/httpx"
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
	err := confx.MustLoad("config.yaml").Parse(cfg)
	if err != nil {
		panic(err)
	}

	// Init logger
	logx.InitLogger(
		logx.WithLevel(cfg.Log.Level),
		logx.WithFile(cfg.Log.File),
	)

	// create gin router
	r := ginx.New()

	// create cron job
	c := cronx.New()
	c.AddFunc("*/1 * * * *", func() {
		logx.Info("cron job run")
	})

	// start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logx.Infof("Server start at %s", addr)
	srv := srvx.New(
		httpx.New(addr, r),
		c,
	)
	srv.Run()
}
