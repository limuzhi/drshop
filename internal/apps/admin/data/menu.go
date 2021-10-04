package data

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/admin/biz"
	"drpshop/internal/apps/admin/vo"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"strings"
)

type menuRepo struct {
	data *Data
	log  *log.Helper
}

func NewMenuRepo(data *Data, logger log.Logger) biz.MenuRepo {
	return &menuRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "admin/data/menu")),
	}
}

//根据用户ID获取用户的权限(可访问)菜单树
func (r *menuRepo) GetUserMenuTreeByUserId(ctx context.Context, uid int64) (*v1.MenuTreeListRes, error) {
	list, err := r.data.sc.MenuTreeListByUserId(ctx, &v1.MenuUserReq{UserId: uid})
	if err != nil {
		return nil, errors.Unauthorized("MenuTreeError", err.Error())
	}
	return list, nil
}

//获取菜单树
func (r *menuRepo) GetMenuTreeList(ctx context.Context) (*v1.MenuTreeListRes, error) {
	return r.data.sc.MenuTreeList(ctx, &v1.MenuListReq{})
}

//根据用户ID获取用户的权限(可访问)菜单列表
func (r *menuRepo) GetUserMenuListByUserId(ctx context.Context, uid int64) (*v1.MenuListRes, error) {
	return r.data.sc.MenuListByUserId(ctx, &v1.MenuUserReq{UserId: uid})
}

//获取菜单树列表
func (r *menuRepo) GetMenuList(ctx context.Context) (*v1.MenuListRes, error) {
	return r.data.sc.MenuList(ctx, &v1.MenuListReq{})
}

//批量删除菜单
func (r *menuRepo) BatchDeleteMenus(ctx context.Context, menuIds []int64) error {
	_, err := r.data.sc.MenuDelete(ctx, &v1.MenuDeleteReq{MenuIds: menuIds})
	return err
}

//创建菜单
func (r *menuRepo) CreateMenu(ctx context.Context, in *vo.CreateMenuReq) error {
	_, err := r.data.sc.MenuAdd(ctx, &v1.MenuAddReq{
		CreateBy:   in.UserId,
		Pid:        in.Pid,
		Name:       strings.TrimSpace(in.Name),
		Title:      strings.TrimSpace(in.Title),
		Icon:       strings.TrimSpace(in.Icon),
		Sort:       int64(in.Sort),
		AlwaysShow: int64(in.AlwaysShow),
		Status:     int64(in.Status),
		Hidden:     int64(in.Hidden),
		Breadcrumb: int64(in.Breadcrumb),
		Path:       strings.TrimSpace(in.Path),
		Component:  strings.TrimSpace(in.Component),
		NoCache:    int64(in.NoCache),
		ActiveMenu: strings.TrimSpace(in.ActiveMenu),
		Redirect:   strings.TrimSpace(in.Redirect),
	})
	return err
}

//更新菜单
func (r *menuRepo) UpdateMneu(ctx context.Context, in *vo.UpdateMenuReq) error {
	_, err := r.data.sc.MenuUpdate(ctx, &v1.MenuUpdateReq{
		UpdateBy:   in.UserId,
		MenuId:     in.MenuId,
		Pid:        in.Pid,
		Name:       strings.TrimSpace(in.Name),
		Title:      strings.TrimSpace(in.Title),
		Icon:       strings.TrimSpace(in.Icon),
		Sort:       int64(in.Sort),
		AlwaysShow: int64(in.AlwaysShow),
		Status:     int64(in.Status),
		Hidden:     int64(in.Hidden),
		Breadcrumb: int64(in.Breadcrumb),
		Path:       strings.TrimSpace(in.Path),
		Component:  strings.TrimSpace(in.Component),
		NoCache:    int64(in.NoCache),
		ActiveMenu: strings.TrimSpace(in.ActiveMenu),
		Redirect:   strings.TrimSpace(in.Redirect),
	})
	return err
}
