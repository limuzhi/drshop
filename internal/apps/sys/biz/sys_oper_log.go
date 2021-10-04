package biz

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/data/model"
	"drpshop/internal/apps/sys/global"
	"github.com/go-kratos/kratos/v2/log"
)

type SysOperLogRepo interface {
	//列表
	ListOperLog(ctx context.Context, req *v1.OperLogListReq) ([]*model.SysOperLog, int64, error)
	//创建
	BatchInsertOperLog(ctx context.Context, ins []*model.SysOperLog) error
	//根据ID获取
	GetInfoById(ctx context.Context, id int64) (*model.SysOperLog, error)
	//批量删除接口
	BatchDeleteByIds(ctx context.Context, ids []int64) error
	//清空
	ClearOperLog(ctx context.Context) error
}

type SysOperLogUsecase struct {
	repo SysOperLogRepo
	log  *log.Helper
}

func NewSysOperLogUsecase(repo SysOperLogRepo, logger log.Logger) *SysOperLogUsecase {
	return &SysOperLogUsecase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "sys/biz/sys_operlog")),
	}
}

//列表
func (uc *SysOperLogUsecase) ListOperLog(ctx context.Context, req *v1.OperLogListReq) ([]*model.SysOperLog, int64, error) {
	if req.PageInfo.PageNum <= 0 {
		req.PageInfo.PageNum = 1
	}
	if req.PageInfo.PageSize <= 0 {
		req.PageInfo.PageSize = 10
	}
	return uc.repo.ListOperLog(ctx, req)
}

//创建
func (uc *SysOperLogUsecase) BatchInsertOperLog(ctx context.Context, ins []*model.SysOperLog) error {
	return uc.repo.BatchInsertOperLog(ctx, ins)
}

//根据ID获取
func (uc *SysOperLogUsecase) GetInfoById(ctx context.Context, id int64) (*model.SysOperLog, error) {
	return uc.repo.GetInfoById(ctx, id)
}

//批量删除接口
func (uc *SysOperLogUsecase) BatchDeleteByIds(ctx context.Context, ids []int64) error {
	return uc.repo.BatchDeleteByIds(ctx, ids)
}

//清空
func (uc *SysOperLogUsecase) ClearOperLog(ctx context.Context) error {
	return uc.repo.ClearOperLog(ctx)
}

//
func (uc *SysOperLogUsecase) DtoOut(v *model.SysOperLog) *v1.OperLogInfoRes {
	info := &v1.OperLogInfoRes{
		OperId:        v.OperId,
		Title:         v.Title,
		BusinessType:  int32(v.BusinessType),
		Method:        v.Method,
		RequestMethod: v.RequestMethod,
		OperatorType:  int32(v.OperatorType),
		OperName:      v.OperName,
		OperUrl:       v.OperUrl,
		OperIp:        v.OperIp,
		OperLocation:  v.OperLocation,
		OperParam:     v.OperParam,
		JsonResult:    v.JsonResult,
		Status:        v.Status,
		ErrorMsg:      v.ErrorMsg,
		OperTime:      global.GetDateByUnix(v.OperTime),
		TimeCost:      v.TimeCost,
		UserId:        v.UserId,
	}
	return info
}
