package biz

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/admin/vo"
	"drpshop/internal/pkg/response"
	"drpshop/pkg/token"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/util/gconv"
	"strings"
)

type RoleRepo interface {
	SearchRoleList(ctx context.Context, in *v1.RoleListReq) (*v1.RoleListRes, error)
	BatchDeleteRole(ctx context.Context, ids []int64) error
	InsertRole(ctx context.Context, in *vo.RoleAddReq) error
	EditRole(ctx context.Context, in *vo.RoleUpdateReq) error
	ChangeRoleStatus(ctx context.Context, id, status int64) error
	RoleMeunsByRoleIds(ctx context.Context, roleId int64) (*v1.QueryMenuByRoleIdRes, error)
	RoleApisByRoleIds(ctx context.Context, roleId int64) (*v1.QueryApisByRoleIdRes, error)

	RoleMenusUpdate(ctx context.Context, userId, roleId int64, menuIds []int64) error
	RoleApisUpdate(ctx context.Context, userId, roleId int64, apiIds []int64) error
}

type RoleController struct {
	roleRepo RoleRepo
	log      *log.Helper
	response.Api
}

func NewRoleController(repo RoleRepo, logger log.Logger) *RoleController {
	log := log.NewHelper(log.With(logger, "module", "admin/biz/role"))
	return &RoleController{
		roleRepo: repo,
		log:      log,
	}
}

func (e *RoleController) RoleList(c *gin.Context) {
	var req vo.RoleSearchReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	reply, err := e.roleRepo.SearchRoleList(c, &v1.RoleListReq{
		Name:     strings.TrimSpace(req.Name),
		RoleKey:  strings.TrimSpace(req.RoleKey),
		Status:   gconv.Int64(req.Status),
		PageNum:  int64(req.GetPageIndex()),
		PageSize: int64(req.GetPageSize()),
	})
	if err != nil {
		e.Api.Fail(c, 40001, err)
		return
	}
	e.Api.PageSuccess(c, reply.List, int(reply.Total), req.PageNum, req.PageSize)
}

func (e *RoleController) RoleCreate(c *gin.Context) {
	var req vo.RoleAddReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	err := e.roleRepo.InsertRole(c, &req)
	if err != nil {
		e.Api.Fail(c, 40002, err)
		return
	}
	e.Api.Success(c, nil)
}

func (e *RoleController) RoleUpdate(c *gin.Context) {
	var req vo.RoleUpdateReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	err := e.roleRepo.EditRole(c, &req)
	if err != nil {
		e.Api.Fail(c, 40002, err)
		return
	}
	e.Api.Success(c, nil)
}

func (e *RoleController) BatchDelete(c *gin.Context) {
	var req vo.RoleBatchDeleteReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	if len(req.RoleIds) == 0 {
		e.Api.Fail(c, 10001, errors.New("参数错误"))
		return
	}
	err := e.roleRepo.BatchDeleteRole(c, req.RoleIds)
	if err != nil {
		e.Api.Fail(c, 40002, err)
		return
	}
	e.Api.Success(c, nil)
}

func (e *RoleController) ChangeStatus(c *gin.Context) {
	var req vo.RoleStatusReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	if req.RoleId == 1 {
		e.Api.Fail(c, 40101, errors.New("超级管理员不能操作"))
		return
	}
	err := e.roleRepo.ChangeRoleStatus(c, req.RoleId, req.Status)
	if err != nil {
		e.Api.Fail(c, 40003, err)
		return
	}
	e.Api.Success(c, nil)
}

func (e *RoleController) RoleMenusList(c *gin.Context) {
	var req vo.RoleMenusApisListReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	reply, err := e.roleRepo.RoleMeunsByRoleIds(c.Request.Context(), req.RoleId)
	if err != nil {
		e.Api.Fail(c, 40002, err)
		return
	}
	e.Api.Success(c, reply.List)
}

func (e *RoleController) RoleApisList(c *gin.Context) {
	var req vo.RoleMenusApisListReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	reply, err := e.roleRepo.RoleApisByRoleIds(c, req.RoleId)
	if err != nil {
		e.Api.Fail(c, 40002, err)
		return
	}
	e.Api.Success(c, reply.List)
}

func (e *RoleController) RoleMenusUpdate(c *gin.Context) {
	var req vo.RoleMenusUpdateReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	userId := token.FormGlobalUidContext(c.Request.Context())
	if userId < 1 {
		e.Api.Fail(c, 10001, errors.New("用户未登录"))
		return
	}
	err := e.roleRepo.RoleMenusUpdate(c, userId, req.RoleId, req.MenuIds)
	if err != nil {
		e.Api.Fail(c, 40002, err)
		return
	}
	e.Api.Success(c, nil)
}

func (e *RoleController) RoleApisUpdate(c *gin.Context) {
	var req vo.RoleApisUpdateReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	userId := e.GetUserId(c)
	if userId < 1 {
		e.Api.Fail(c, 10001, errors.New("用户未登录"))
		return
	}
	err := e.roleRepo.RoleApisUpdate(e.NewContext(c), userId, req.RoleId, req.ApiIds)
	if err != nil {
		e.Api.Fail(c, 40002, err)
		return
	}
	e.Api.Success(c, nil)
}
