package cron

import (
	"context"

	"github.com/nurcahyaari/ecommerce/src/domain/service"
	"github.com/robfig/cron/v3"
)

type CronHandle struct {
	orderService service.OrderServicer
}

func NewCronhandler(
	orderService service.OrderServicer,
) *CronHandle {
	return &CronHandle{
		orderService: orderService,
	}
}

func (c *CronHandle) Router(cron *cron.Cron) {
	cron.AddFunc("* * * * *", func() {
		c.orderService.ExpiredOrder(context.TODO())
		return
	})
}
