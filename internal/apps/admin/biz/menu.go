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
)

type MenuRepo interface {
	//根据用户ID获取用户的权限(可访问)菜单树
	GetUserMenuTreeByUserId(ctx context.Context, uid int64) (*v1.MenuTreeListRes, error)
	//获取菜单树
	GetMenuTreeList(ctx context.Context) (*v1.MenuTreeListRes, error)
	//根据用户ID获取用户的权限(可访问)菜单列表
	GetUserMenuListByUserId(ctx context.Context, uid int64) (*v1.MenuListRes, error)
	//获取菜单树列表
	GetMenuList(ctx context.Context) (*v1.MenuListRes, error)
	//批量删除菜单
	BatchDeleteMenus(ctx context.Context, menuIds []int64) error
	//创建菜单
	CreateMenu(ctx context.Context, in *vo.CreateMenuReq) error
	//更新菜单
	UpdateMneu(ctx context.Context, in *vo.UpdateMenuReq) error
}

type MenuController struct {
	menuRepo MenuRepo
	log      *log.Helper
	response.Api
}

func NewMenuController(repo MenuRepo, logger log.Logger) *MenuController {
	log := log.NewHelper(log.With(logger, "module", "admin/biz/menu"))
	return &MenuController{
		menuRepo: repo,
		log:      log,
	}
}

//获取菜单树
func (e *MenuController) GetMenuTree(c *gin.Context) {
	reply, err := e.menuRepo.GetMenuTreeList(e.NewContext(c))
	if err != nil {
		e.Api.Fail(c, 20003, err)
		return
	}
	e.Api.Success(c, reply.List)
}

//获取菜单列表
func (e *MenuController) GetMenuList(c *gin.Context) {
	reply, err := e.menuRepo.GetMenuList(e.NewContext(c))
	if err != nil {
		e.Api.Fail(c, 10002, err)
		return
	}
	e.Api.Success(c, reply.List)
}

//创建菜单
func (e *MenuController) MenuCreate(c *gin.Context) {
	var req vo.CreateMenuReq
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
	req.UserId = userInfo.UserID
	if err := e.menuRepo.CreateMenu(e.NewContext(c), &req); err != nil {
		e.Api.Fail(c, 10002, err)
		return
	}
	e.Api.Success(c, nil)
}

//更新菜单
func (e *MenuController) MenuUpdate(c *gin.Context) {
	var req vo.UpdateMenuReq
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
	req.UserId = userInfo.UserID
	if err := e.menuRepo.UpdateMneu(e.NewContext(c), &req); err != nil {
		e.Api.Fail(c, 10002, err)
		return
	}
	e.Api.Success(c, nil)
}

//删除菜单
func (e *MenuController) MenubatchDelete(c *gin.Context) {
	var req vo.DeleteMenuReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	if len(req.MenuIds) == 0 {
		e.Api.Fail(c, 10001, errors.New("menuIds参数值为空"))
		return
	}
	userInfo := e.GetUserInfo(c)
	if userInfo == nil || userInfo.UserID <= 0 {
		e.Api.Fail(c, 20001, errors.New("未登录"))
		return
	}
	if err := e.menuRepo.BatchDeleteMenus(e.NewContext(c), req.MenuIds); err != nil {
		if err != nil {
			e.Api.Fail(c, 20002, err)
			return
		}
	}
	e.Api.Success(c, nil)
}

//获取用户的可访问菜单列表
func (e *MenuController) GetUserMenuListByUserId(c *gin.Context) {
	userInfo := token.FormLoginContext(c.Request.Context())
	if userInfo == nil || userInfo.UserID <= 0 {
		e.Api.Fail(c, 20001, errors.New("未登录"))
		return
	}
	reply, err := e.menuRepo.GetUserMenuListByUserId(e.NewContext(c), userInfo.UserID)
	if err != nil {
		e.Api.Fail(c, 20002, err)
		return
	}
	e.Api.Success(c, reply.List)
}

//获取用户的可访问菜单树
func (e *MenuController) GetUserMenuTreeByUserId(c *gin.Context) {
	userInfo := token.FormLoginContext(c.Request.Context())
	if userInfo == nil || userInfo.UserID <= 0 {
		e.Api.Fail(c, 20001, errors.New("未登录"))
		return
	}
	res, err := e.menuRepo.GetUserMenuTreeByUserId(e.NewContext(c), userInfo.UserID)
	if err != nil {
		e.Api.Fail(c, 20002, err)
		return
	}
	e.Api.Success(c, res.List)
}
