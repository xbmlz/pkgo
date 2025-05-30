package cronx

import (
	"github.com/robfig/cron/v3"
	"github.com/xbmlz/pkgo/srvx"
)

var _ srvx.Server = (*Cronx)(nil)

type Cronx struct {
	*cron.Cron
}

func New() *Cronx {
	return &Cronx{cron.New()}
}

func (c *Cronx) Run() error {
	c.Start()
	return nil
}

func (c *Cronx) Shutdown() error {
	c.Stop()
	return nil
}
