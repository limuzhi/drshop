package data

import (
	"drpshop/internal/apps/admin/biz"
	"github.com/go-kratos/kratos/v2/log"
)

type monitorRepo struct {
	data *Data
	log  *log.Helper
}

func NewMonitorRepo(data *Data, logger log.Logger) biz.MonitorRepo {
	return &monitorRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "admin/data/monitor")),
	}
}
