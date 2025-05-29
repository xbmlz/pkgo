package main

import (
	"fmt"

	"github.com/xbmlz/pkgo/config"
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
	var c Config
	config.MustLoad("config.yaml").MustParse(&c)

	// Init logger
	log.InitLogger(
		log.WithLevel(c.Log.Level),
		log.WithFile(c.Log.File),
	)

	// create gin router
	r := ginx.New()

	// start server
	addr := fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
	log.Infof("Server start at %s", addr)
	server.NewHTTPServer(r, addr).Run()
}
