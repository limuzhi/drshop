package biz

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/admin/vo"
	"drpshop/internal/pkg/response"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"strings"
)

type DictDataRepo interface {
	GetDictDataListSelect(ctx context.Context, dictType string) (*v1.DictDataListRes, error)
	DictDataList(ctx context.Context, in *vo.DictDataListReq) (*v1.DictDataListRes, error)
	DictDataDetail(ctx context.Context, id int64) (*v1.DictDataListData, error)
	DictDataCreate(ctx context.Context, in *vo.DictDataCreateReq) error
	DictDataUpdate(ctx context.Context, in *vo.DictDataUpdateReq) error
	DictDataDelete(ctx context.Context, ids []int64) error
}

type DictDataController struct {
	dictDataRepo DictDataRepo
	log          *log.Helper
	response.Api
}

func NewDictDataController(repo DictDataRepo, logger log.Logger) *DictDataController {
	log := log.NewHelper(log.With(logger, "module", "admin/biz/dict_data"))
	return &DictDataController{
		dictDataRepo: repo,
		log:          log,
	}
}

func (e *DictDataController) DictDataSelect(c *gin.Context) {
	var req vo.DictDataSelectReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	dictType := strings.TrimSpace(req.DictType)
	if dictType == "" {
		e.Api.Fail(c, 10001, errors.New("dictType不能为空"))
		return
	}
	res, err := e.dictDataRepo.GetDictDataListSelect(e.NewContext(c), dictType)
	if err != nil {
		e.Api.Fail(c, 30002, errors.New("获取失败"))
		return
	}
	e.Api.Success(c, res.List)
}

func (e *DictDataController) DictDataList(c *gin.Context) {
	var req vo.DictDataListReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	reply, err := e.dictDataRepo.DictDataList(e.NewContext(c), &req)
	if err != nil {
		e.Api.Fail(c, 50001, err)
		return
	}
	e.Api.PageSuccess(c, reply.List, int(reply.Total), req.GetPageIndex(), req.GetPageSize())
}

func (e *DictDataController) DictDataDetail(c *gin.Context) {
	var req vo.DictDataDetailReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	reply, err := e.dictDataRepo.DictDataDetail(e.NewContext(c), req.DictCode)
	if err != nil {
		e.Api.Fail(c, 60001, err)
		return
	}
	e.Api.Success(c, reply)
}
func (e *DictDataController) DictDataCreate(c *gin.Context) {
	var req vo.DictDataCreateReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	userInfo := e.GetUserInfo(c)
	if userInfo == nil || userInfo.UserID <= 0 {
		e.Api.Fail(c, 20001, errors.New("未登录"))
		return
	}
	req.CreateBy = userInfo.UserID
	err := e.dictDataRepo.DictDataCreate(e.NewContext(c), &req)
	if err != nil {
		e.Api.Fail(c, 60001, err)
		return
	}
	e.Api.Success(c, nil)
}
func (e *DictDataController) DictDataUpdate(c *gin.Context) {
	var req vo.DictDataUpdateReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	userInfo := e.GetUserInfo(c)
	if userInfo == nil || userInfo.UserID <= 0 {
		e.Api.Fail(c, 20001, errors.New("未登录"))
		return
	}
	req.UpdateBy = userInfo.UserID
	err := e.dictDataRepo.DictDataUpdate(e.NewContext(c), &req)
	if err != nil {
		e.Api.Fail(c, 60001, err)
		return
	}
	e.Api.Success(c, nil)
}
func (e *DictDataController) DictDataDelete(c *gin.Context) {
	var req vo.DictDataDeleteReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	if len(req.DictCodes) == 0 {
		e.Api.Fail(c, 20001, errors.New("参数不能为空"))
		return
	}
	err := e.dictDataRepo.DictDataDelete(e.NewContext(c), req.DictCodes)
	if err != nil {
		e.Api.Fail(c, 60001, err)
		return
	}
	e.Api.Success(c, nil)
}
