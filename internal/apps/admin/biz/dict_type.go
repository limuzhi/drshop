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

type DictTypeRepo interface {
	DictTypeList(ctx context.Context, in *vo.DictTypeListReq) (*v1.DictTypeListRes, error)
	DictTypeSelect(ctx context.Context) (*v1.DictTypeListRes, error)
	DictTypeDetail(ctx context.Context, id int64) (*v1.DictTypeListData, error)
	DictTypeCreate(ctx context.Context, in *vo.DictTypeCreateReq) error
	DictTypeUpdate(ctx context.Context, in *vo.DictTypeUpdateReq) error
	DictTypeDelete(ctx context.Context, ids []int64) error
}

type DictTypeController struct {
	dictTypeRepo DictTypeRepo
	log          *log.Helper
	response.Api
}

func NewDictTypeController(repo DictTypeRepo, logger log.Logger) *DictTypeController {
	log := log.NewHelper(log.With(logger, "module", "admin/biz/dict_type"))
	return &DictTypeController{
		dictTypeRepo: repo,
		log:          log,
	}
}

func (e *DictTypeController) DictTypeSelect(c *gin.Context) {
	reply, err := e.dictTypeRepo.DictTypeSelect(c)
	if err != nil {
		e.Api.Fail(c, 60001, err)
		return
	}
	e.Api.Success(c, reply.List)
}

func (e *DictTypeController) DictTypeList(c *gin.Context) {
	var req vo.DictTypeListReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	reply, err := e.dictTypeRepo.DictTypeList(e.NewContext(c), &req)
	if err != nil {
		e.Api.Fail(c, 50001, err)
		return
	}
	e.Api.PageSuccess(c, reply.List, int(reply.Total), req.GetPageIndex(), req.GetPageSize())
}

func (e *DictTypeController) DictTypeDetail(c *gin.Context) {
	var req vo.DictTypeDetailReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	reply, err := e.dictTypeRepo.DictTypeDetail(e.NewContext(c), req.DictId)
	if err != nil {
		e.Api.Fail(c, 60001, err)
		return
	}
	e.Api.Success(c, reply)
}

func (e *DictTypeController) DictTypeCreate(c *gin.Context) {
	var req vo.DictTypeCreateReq
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
	err := e.dictTypeRepo.DictTypeCreate(e.NewContext(c), &req)
	if err != nil {
		e.Api.Fail(c, 60001, err)
		return
	}
	e.Api.Success(c, nil)
}
func (e *DictTypeController) DictTypeUpdate(c *gin.Context) {
	var req vo.DictTypeUpdateReq
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
	err := e.dictTypeRepo.DictTypeUpdate(e.NewContext(c), &req)
	if err != nil {
		e.Api.Fail(c, 60001, err)
		return
	}
	e.Api.Success(c, nil)
}
func (e *DictTypeController) DictTypeDelete(c *gin.Context) {
	var req vo.DictTypeDeleteReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	if len(req.DictIds) == 0 {
		e.Api.Fail(c, 20001, errors.New("参数不能为空"))
		return
	}
	err := e.dictTypeRepo.DictTypeDelete(e.NewContext(c), req.DictIds)
	if err != nil {
		e.Api.Fail(c, 60001, err)
		return
	}
	e.Api.Success(c, nil)
}
