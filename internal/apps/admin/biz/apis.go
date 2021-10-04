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

type ApisRepo interface {
	GetApiTreeList(ctx context.Context) (*v1.ApisTreeListRes, error)
	GetApiList(ctx context.Context, in *vo.ApiListReq) (*v1.ApisListRes, error)
	ApiCreate(ctx context.Context, in *vo.CreateApiReq) error
	ApiUpdate(ctx context.Context, in *vo.UpdateApiReq) error
	ApiBatchDelete(ctx context.Context, ids []int64) error
}

type ApisController struct {
	apisRepo ApisRepo
	log      *log.Helper
	response.Api
}

func NewApisController(repo ApisRepo, logger log.Logger) *ApisController {
	log := log.NewHelper(log.With(logger, "module", "admin/biz/apis"))
	return &ApisController{
		apisRepo: repo,
		log:      log,
	}
}

func (e *ApisController) GetApiTree(c *gin.Context) {
	reply, err := e.apisRepo.GetApiTreeList(e.NewContext(c))
	if err != nil {
		e.Api.Fail(c, 50001, err)
		return
	}
	e.Api.Success(c, reply.List)
}

func (e *ApisController) GetApiList(c *gin.Context) {
	var req vo.ApiListReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	reply, err := e.apisRepo.GetApiList(e.NewContext(c), &req)
	if err != nil {
		e.Api.Fail(c, 50001, err)
		return
	}
	e.Api.PageSuccess(c, reply.List, int(reply.Total),
		req.GetPageIndex(), req.GetPageSize())
}

func (e *ApisController) ApiCreate(c *gin.Context) {
	var req vo.CreateApiReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	req.CreateBy = e.GetUserId(c)
	err := e.apisRepo.ApiCreate(e.NewContext(c), &req)
	if err != nil {
		e.Api.Fail(c, 50002, err)
		return
	}
	e.Api.Success(c, nil)
}

func (e *ApisController) ApiUpdate(c *gin.Context) {
	var req vo.UpdateApiReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	req.UpdateBy = e.GetUserId(c)
	err := e.apisRepo.ApiUpdate(e.NewContext(c), &req)
	if err != nil {
		e.Api.Fail(c, 50002, err)
		return
	}
	e.Api.Success(c, nil)
}

func (e *ApisController) ApiBatchDelete(c *gin.Context) {
	var req vo.BatchDeleteApiReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	if len(req.ApiIds) == 0 {
		e.Api.Fail(c, 10001, errors.New("参数不能为空"))
		return
	}
	err := e.apisRepo.ApiBatchDelete(e.NewContext(c), req.ApiIds)
	if err != nil {
		e.Api.Fail(c, 50003, err)
		return
	}
	e.Api.Success(c, nil)
}
