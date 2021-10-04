package biz

import (
	"context"
	"drpshop/internal/pkg/response"
	"github.com/go-kratos/kratos/v2/log"
)

type CasbinRepo interface {
	CheckCasbin(ctx context.Context, userId int64, obj, act string) bool
}

type CasbinController struct {
	casbinRepo CasbinRepo
	log    *log.Helper
	response.Api
}

func NewCasbinController(repo CasbinRepo, logger log.Logger) *CasbinController {
	log := log.NewHelper(log.With(logger, "module", "admin/biz/casbin"))
	return &CasbinController{
		casbinRepo: repo,
		log:    log,
	}
}

func (e *CasbinController) CheckAuth(ctx context.Context,userId int64, obj, act string) bool {
	return e.casbinRepo.CheckCasbin(ctx,userId,obj,act)
}
