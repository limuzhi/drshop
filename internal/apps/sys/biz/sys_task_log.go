package biz

import "github.com/go-kratos/kratos/v2/log"

type SysTaskLogRepo interface {
}
type SysTaskLogUsecase struct {
	repo SysTaskLogRepo
	log  *log.Helper
}

func NewSysTaskLogUsecase(repo SysTaskLogRepo, logger log.Logger) *SysTaskLogUsecase {
	return &SysTaskLogUsecase{repo: repo, log: log.NewHelper(
		log.With(logger, "module", "sys/biz/sys_task_log"))}
}
