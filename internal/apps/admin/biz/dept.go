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

type DeptRepo interface {
	DeptList(ctx context.Context, in *vo.DeptListReq) (*v1.DeptListRes, error)
	DeptTree(ctx context.Context, in *vo.DeptListReq) (*v1.DeptListRes, error)
	DeptDetail(ctx context.Context, deptId int64) (*v1.DeptDetailRes, error)
	DeptCreate(ctx context.Context, in *vo.DeptCreateReq) error
	DeptUpdate(ctx context.Context, in *vo.DeptUpdateReq) error
	DeptDelete(ctx context.Context, ids []int64) error
}

type DeptController struct {
	deptRepo DeptRepo
	log      *log.Helper
	response.Api
}

func NewDeptController(repo DeptRepo, logger log.Logger) *DeptController {
	log := log.NewHelper(log.With(logger, "module", "admin/biz/dept"))
	return &DeptController{
		deptRepo: repo,
		log:      log,
	}
}

func (e *DeptController) DeptList(c *gin.Context) {
	var req vo.DeptListReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	reply, err := e.deptRepo.DeptList(e.NewContext(c), &req)
	if err != nil {
		e.Api.Fail(c, 60001, err)
		return
	}
	e.Api.Success(c, reply.List)
}

func (e *DeptController) DeptTree(c *gin.Context) {
	var req vo.DeptListReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	reply, err := e.deptRepo.DeptTree(e.NewContext(c), &req)
	if err != nil {
		e.Api.Fail(c, 60001, err)
		return
	}
	e.Api.Success(c, reply.List)
}

func (e *DeptController) DeptDetail(c *gin.Context) {
	var req vo.DeptDetailReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	reply, err := e.deptRepo.DeptDetail(e.NewContext(c), req.DeptId)
	if err != nil {
		e.Api.Fail(c, 60001, err)
		return
	}
	e.Api.Success(c, reply)
}
func (e *DeptController) DeptCreate(c *gin.Context) {
	var req vo.DeptCreateReq
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
	err := e.deptRepo.DeptCreate(e.NewContext(c), &req)
	if err != nil {
		e.Api.Fail(c, 60001, err)
		return
	}
	e.Api.Success(c, nil)
}

func (e *DeptController) DeptUpdate(c *gin.Context) {
	var req vo.DeptUpdateReq
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
	err := e.deptRepo.DeptUpdate(e.NewContext(c), &req)
	if err != nil {
		e.Api.Fail(c, 60001, err)
		return
	}
	e.Api.Success(c, nil)
}
func (e *DeptController) DeptDelete(c *gin.Context) {
	var req vo.DeptDeleteReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	if len(req.DeptIds) == 0 {
		e.Api.Fail(c, 20001, errors.New("参数不能为空"))
		return
	}
	err := e.deptRepo.DeptDelete(e.NewContext(c), req.DeptIds)
	if err != nil {
		e.Api.Fail(c, 60001, err)
		return
	}
	e.Api.Success(c, nil)
}
