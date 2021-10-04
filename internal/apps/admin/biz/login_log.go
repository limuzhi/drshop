package biz

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/admin/vo"
	"drpshop/internal/pkg/response"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
)

type LoginLogRepo interface {
	List(ctx context.Context, in *vo.LoginLogSearchReq) (*v1.LoginlogListRes, error)
	BatchDelete(ctx context.Context, ids []int64) error
	Clear(ctx context.Context) error
}

type LoginLogController struct {
	loginLogRepo LoginLogRepo
	log          *log.Helper
	response.Api
}

func NewLoginLogController(repo LoginLogRepo, logger log.Logger) *LoginLogController {
	log := log.NewHelper(log.With(logger, "module", "admin/biz/login_log"))
	return &LoginLogController{
		loginLogRepo: repo,
		log:          log,
	}
}

func (uc *LoginLogController) LoginList(c *gin.Context) {
	var req vo.LoginLogSearchReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		uc.Api.Fail(c, 10001, err)
		return
	}
	list, err := uc.loginLogRepo.List(uc.NewContext(c), &req)
	if err != nil {
		uc.Api.Fail(c, 30001, err)
		return
	}
	uc.Api.PageSuccess(c, list.List, int(list.Total), req.GetPageIndex(), req.GetPageSize())
}

func (uc *LoginLogController) BatchDeleteLog(c *gin.Context) {
	var req vo.LoginLogDeleteReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		uc.Api.Fail(c, 10001, err)
		return
	}
	if len(req.LoginIds) == 0 {
		uc.Api.Fail(c, 30001, errors.New("参数值为空"))
		return
	}
	err := uc.loginLogRepo.BatchDelete(uc.NewContext(c), req.LoginIds)
	if err != nil {
		uc.Api.Fail(c, 30001, err)
		return
	}
	uc.Api.Success(c, nil)
}

func (uc *LoginLogController) Clearlog(c *gin.Context) {
	err := uc.loginLogRepo.Clear(uc.NewContext(c))
	if err != nil {
		uc.Api.Fail(c, 30001, err)
		return
	}
	uc.Api.Success(c, nil)
}
