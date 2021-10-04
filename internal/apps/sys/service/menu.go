package service

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/data/model"
	"github.com/go-kratos/kratos/v2/errors"
)

//获取菜单列表
func (s *SysService) MenuList(ctx context.Context, _ *v1.MenuListReq) (*v1.MenuListRes, error) {
	reply, err := s.menuc.ListMenu(ctx)
	if err != nil {
		return nil, err
	}
	list := s.menuc.GetMenusListDto(reply)
	return &v1.MenuListRes{List: list, Total: int64(len(list))}, nil
}

//获取菜单树
func (s *SysService) MenuTreeList(ctx context.Context, _ *v1.MenuListReq) (*v1.MenuTreeListRes, error) {
	list, err := s.menuc.ListMenuTree(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.MenuTreeListRes{List: list}, nil
}

//创建菜单
func (s *SysService) MenuAdd(ctx context.Context, in *v1.MenuAddReq) (*v1.CommonRes, error) {
	data := &model.SysMenu{
		Pid:        in.Pid,
		Name:       in.Name,
		Title:      in.Title,
		Icon:       in.Icon,
		Sort:       int(in.Sort),
		Hidden:     int(in.Hidden),
		NoCache:    int(in.NoCache),
		Path:       in.Path,
		Redirect:   in.Redirect,
		ActiveMenu: in.ActiveMenu,
		JumpPath:   in.JumpPath,
		Component:  in.Component,
		IsFrame:    int(in.IsFrame),
		ModuleType: in.ModuleType,
		ModelId:    int(in.ModelId),
		Status:     int(in.Status),
		AlwaysShow: int(in.AlwaysShow),
		Breadcrumb: int(in.Breadcrumb),
	}
	data.CreateBy = in.CreateBy
	err := s.menuc.CreateMenu(ctx, data)
	return &v1.CommonRes{}, err
}

//更新菜单
func (s *SysService) MenuUpdate(ctx context.Context, in *v1.MenuUpdateReq) (*v1.CommonRes, error) {
	data := &model.SysMenu{
		MenuId:     in.MenuId,
		Pid:        in.Pid,
		Name:       in.Name,
		Title:      in.Title,
		Icon:       in.Icon,
		Sort:       int(in.Sort),
		Status:     int(in.Status),
		Hidden:     int(in.Hidden),
		Path:       in.Path,
		JumpPath:   in.JumpPath,
		Component:  in.Component,
		IsFrame:    int(in.IsFrame),
		ModuleType: in.ModuleType,
		ModelId:    int(in.ModelId),
		NoCache:    int(in.NoCache),
		ActiveMenu: in.ActiveMenu,
		Redirect:   in.Redirect,
		AlwaysShow: int(in.AlwaysShow),
		Breadcrumb: int(in.Breadcrumb),
	}
	data.UpdateBy = in.UpdateBy
	err := s.menuc.UpdateMenu(ctx, in.MenuId, data)
	return &v1.CommonRes{}, err
}

//批量删除菜单
func (s *SysService) MenuDelete(ctx context.Context, in *v1.MenuDeleteReq) (*v1.CommonRes, error) {
	err := s.menuc.BatchDeleteByIds(ctx, in.MenuIds)
	return &v1.CommonRes{}, err
}

//获取用户的可访问菜单列表
func (s *SysService) MenuListByUserId(ctx context.Context, in *v1.MenuUserReq) (*v1.MenuListRes, error) {
	reply, err := s.menuc.GetUserMenusByUserId(ctx, in.UserId)
	if err != nil {
		return nil, errors.Unauthorized("MenuTreeListByUserIdError", "获取用户的可访问菜单列表--错误")
	}
	list := s.menuc.GetMenusListDto(reply)
	return &v1.MenuListRes{List: list, Total: int64(len(list))}, nil
}

func (s *SysService) MenuTreeListByUserId(ctx context.Context, req *v1.MenuUserReq) (*v1.MenuTreeListRes, error) {
	list, err := s.menuc.GetUserMenuTreeByUserId(ctx, req.UserId)
	if err != nil {
		return nil, errors.Unauthorized("MenuTreeListByUserIdError", "用户的权限(可访问)菜单列表错误")
	}
	return &v1.MenuTreeListRes{List: list}, nil
}
