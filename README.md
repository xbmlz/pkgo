# pkgo

Package pkgo is a collection of packages for Go.

<a title="Build Status" target="_blank" href="https://github.com/xbmlz/pkgo/actions/workflows/test.yml"><img src="https://img.shields.io/github/actions/workflow/status/xbmlz/pkgo/test.yml?style=flat-square"></a>
<a title="GoDoc" target="_blank" href="https://godoc.org/github.com/xbmlz/pkgo"><img src="http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square"></a>
<a title="Go Report Card" target="_blank" href="https://goreportcard.com/report/github.com/xbmlz/pkgo"><img src="https://goreportcard.com/badge/github.com/xbmlz/pkgo?style=flat-square"></a>
<a title="Coverage Status" target="_blank" href="https://coveralls.io/github/xbmlz/pkgo"><img src="https://img.shields.io/coveralls/github/xbmlz/pkgo.svg?style=flat-square&color=CC9933"></a>
<a title="Code Size" target="_blank" href="https://github.com/xbmlz/pkgo"><img src="https://img.shields.io/github/languages/code-size/xbmlz/pkgo.svg?style=flat-square"></a>

## Get started

```go
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
	log.InitLogger()

	// create gin router
	r := ginx.New()

	// start server
	addr := fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
	log.Infof("Server start at %s", addr)
	server.NewHTTPServer(r, addr).Run()
}
```