package cron

import (
	"context"

	"github.com/nurcahyaari/ecommerce/config"
	cronhandler "github.com/nurcahyaari/ecommerce/src/handlers/cron"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

type CronScheduler struct {
	cfg     config.Config
	log     zerolog.Logger
	cron    *cron.Cron
	handler *cronhandler.CronHandle
}

func New(
	cfg config.Config,
	handler *cronhandler.CronHandle,
	log zerolog.Logger,
) *CronScheduler {
	return &CronScheduler{
		cfg:     cfg,
		cron:    cron.New(),
		handler: handler,
		log:     log,
	}
}

func (c *CronScheduler) Router() {
	c.handler.Router(c.cron)
}

func (c *CronScheduler) Listen() {
	c.Router()
	c.log.Info().Msg("running cron job")
	c.cron.Start()
}

func (c *CronScheduler) Shutdown(ctx context.Context) error {
	c.cron.Stop()
	return nil
}
