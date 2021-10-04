package biz

import (
	"context"
	"drpshop/api/common"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/admin/vo"
	"drpshop/internal/pkg/response"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"strings"
)

type OperLogRepo interface {
	GetList(ctx context.Context, req *v1.OperLogListReq) (*v1.OperLogListRes, error)
	GetDetail(ctx context.Context, id int64) (*v1.OperLogInfoRes, error)
	BatchDelete(ctx context.Context, ids []int64) error
	Clear(ctx context.Context) error
	SaveOperlogChannel(out <-chan *v1.OperLogSaveData)
}

type OperLogController struct {
	operLogRepo OperLogRepo
	log         *log.Helper
	response.Api
}

func NewOperLogController(repo OperLogRepo, logger log.Logger) *OperLogController {
	log := log.NewHelper(log.With(logger, "module", "admin/biz/oper_log"))
	return &OperLogController{
		operLogRepo: repo,
		log:         log,
	}
}

func (uc *OperLogController) OperLogList(c *gin.Context) {
	var req vo.OperLogSearchReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		uc.Api.Fail(c, 10001, err)
		return
	}
	list, err := uc.operLogRepo.GetList(uc.NewContext(c), &v1.OperLogListReq{
		Title:        strings.TrimSpace(req.Title),
		OperName:     strings.TrimSpace(req.OperName),
		Status:       strings.TrimSpace(req.Status),
		BusinessType: strings.TrimSpace(req.BusinessType),
		PageInfo: &common.PageReq{
			PageSize:  int64(req.GetPageSize()),
			PageNum:   int64(req.GetPageIndex()),
			BeginTime: req.GetBeginTime(),
			EndTime:   req.GetEndTime(),
		},
	})
	if err != nil {
		uc.Api.Fail(c, 30001, err)
		return
	}
	uc.Api.PageSuccess(c, list.List, int(list.Total), req.GetPageIndex(), req.GetPageSize())
}

func (uc *OperLogController) OperLogDetail(c *gin.Context) {
	var req vo.OperLogDetailReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		uc.Api.Fail(c, 10001, err)
		return
	}
	info, err := uc.operLogRepo.GetDetail(uc.NewContext(c), req.OperId)
	if err != nil {
		uc.Api.Fail(c, 30002, err)
	}
	uc.Api.Success(c, info)
}

func (uc *OperLogController) OperLogDelete(c *gin.Context) {
	var req vo.OperLogDeleteReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		uc.Api.Fail(c, 10001, err)
		return
	}
	if len(req.OperIds) == 0 {
		uc.Api.Fail(c, 10001, errors.New("参数未携带"))
		return
	}
	err := uc.operLogRepo.BatchDelete(uc.NewContext(c), req.OperIds)
	if err != nil {
		uc.Api.Fail(c, 30003, errors.New("删除失败"))
		return
	}
	uc.Api.Success(c, nil)
}

func (uc *OperLogController) OperLogClear(c *gin.Context) {
	err := uc.operLogRepo.Clear(uc.NewContext(c))
	if err != nil {
		uc.Api.Fail(c, 30003, errors.New("清空失败"))
		return
	}
	uc.Api.Success(c, nil)
}
