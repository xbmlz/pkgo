package main

import (
	"fmt"

	"github.com/xbmlz/pkgo/conf"
	"github.com/xbmlz/pkgo/server"
	"github.com/xbmlz/pkgo/utils"
)

type Config struct {
	Server server.Config `yaml:"server"`
}

func main() {
	// Load config
	cfg := &Config{}
	err := conf.Load(cfg)
	if err != nil {
		panic(err)
	}

	fmt.Println(utils.OrElse(cfg.Server.ReadTimeout, 10))

	// r := ginx.New()

	// r.GET("/", func(c *gin.Context) {
	// 	ginx.ResponseOk(c, "Hello, World!")
	// })

	// server.Run(cfg.Server, r)
}
