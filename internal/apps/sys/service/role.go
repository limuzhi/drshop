package service

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/global"
	"drpshop/internal/apps/sys/model"
	"errors"
	"fmt"
	"github.com/thoas/go-funk"
)

//创建角色
func (s *SysService) RoleAdd(ctx context.Context, req *v1.RoleAddReq) (*v1.CommonRes, error) {
	err := s.rolec.CreateRole(ctx, &model.SysRole{
		Status:    req.Status,
		Sort:      req.Sort,
		PID:       req.Pid,
		Name:      req.Name,
		RoleKey:   req.RoleKey,
		Remark:    req.Remark,
		DataScope: int32(req.DataScope),
	})
	return &v1.CommonRes{}, err
}

//角色列表数据
func (s *SysService) RoleList(ctx context.Context, in *v1.RoleListReq) (*v1.RoleListRes, error) {
	if in.PageNum <= 0 {
		in.PageNum = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 10
	}
	result, total, err := s.rolec.ListRole(ctx, in)
	if err != nil {
		return nil, err
	}
	list := make([]*v1.RoleListData, len(result))
	for k, v := range result {
		info := &v1.RoleListData{
			RoleId:    v.RoleId,
			Name:      v.Name,
			Pid:       v.PID,
			RoleKey:   v.RoleKey,
			Remark:    v.Remark,
			Sort:      v.Sort,
			DataScope: int64(v.DataScope),
			Status:    v.Status,
		}
		list[k] = info
	}
	return &v1.RoleListRes{Total: total, List: list}, nil
}

//修改用户角色
func (s *SysService) RoleUpdate(ctx context.Context, in *v1.RoleUpdateReq) (*v1.CommonRes, error) {
	if in.RoleId <= 0 {
		return nil, errors.New("请求参数roleId不能小于1")
	}
	err := s.rolec.UpdateRole(ctx, &model.SysRole{
		RoleId:    in.RoleId,
		Name:      in.Name,
		PID:       in.Pid,
		RoleKey:   in.RoleKey,
		Remark:    in.Remark,
		Sort:      in.Sort,
		DataScope: int32(in.DataScope),
		Status:    in.Status,
	})
	return &v1.CommonRes{}, err
}

//删除用户角色
func (s *SysService) RoleDelete(ctx context.Context, in *v1.RoleDeleteReq) (*v1.CommonRes, error) {
	if len(in.RoleIds) == 0 {
		return nil, errors.New("请求参数错误")
	}
	err := s.rolec.BatchDeleteRoleByIds(ctx, in.RoleIds)
	return &v1.CommonRes{}, err
}

//修改角色状态
func (s *SysService) UpdateRoleStatus(ctx context.Context, in *v1.UpdateRoleStatusReq) (*v1.CommonRes, error) {
	err := s.rolec.ChangeRoleStatus(ctx, in.RoleId, in.Status)
	return &v1.CommonRes{}, err
}

//获取角色的权限菜单
func (s *SysService) GetMenusByRoleId(ctx context.Context, in *v1.QueryMenuByRoleIdReq) (*v1.QueryMenuByRoleIdRes, error) {
	reply, err := s.rolec.GetRoleMenuList(ctx, []int64{in.RoleId})
	if err != nil {
		return nil, err
	}
	menuList := make([]*model.SysMenu, 0)
	menuIdMap := make(map[int64]int64)
	for _, v := range reply {
		for _, val := range v.Menus {
			if _, ok := menuIdMap[val.MenuId]; !ok {
				menuList = append(menuList, val)
			}
		}
	}
	list := s.menuc.GetMenusListDto(menuList)
	return &v1.QueryMenuByRoleIdRes{List: list, Total: int64(len(list))}, nil
}

//更新角色的权限菜单
func (s *SysService) UpdateMenuRole(ctx context.Context, in *v1.UpdateMenuRoleReq) (*v1.CommonRes, error) {
	list, err := s.rolec.GetInfo(ctx, []int64{in.RoleId})
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.New("未找到角色信息")
	}
	//TODO
	// 当前用户角色排序最小值（最高等级角色）以及当前用户

	// (非管理员)不能更新比自己角色等级高或相等角色的权限菜单

	ctxUserMenus, err := s.menuc.GetUserMenusByUserId(ctx, in.UserId)
	if err != nil {
		return nil, errors.New("获取当前用户的可访问菜单列表失败")
	}
	// 获取当前用户所拥有的权限菜单ID
	ctxUserMenusIds := make([]int64, 0)
	for _, menu := range ctxUserMenus {
		ctxUserMenusIds = append(ctxUserMenusIds, menu.MenuId)
	}
	// 前端传来最新的MenuIds集合
	menuIds := in.MenuIds

	// 用户需要修改的菜单集合
	reqMenus := make([]*model.SysMenu, 0)
	// (非管理员)不能把角色的权限菜单设置的比当前用户所拥有的权限菜单多
	//
	if in.UserId != 1 {
		for _, id := range menuIds {
			if !funk.Contains(ctxUserMenusIds, id) {
				return nil, errors.New(fmt.Sprintf("无权设置ID为%d的菜单", id))
			}
		}
		for _, id := range menuIds {
			for _, menu := range ctxUserMenus {
				if id == menu.MenuId {
					reqMenus = append(reqMenus, menu)
					break
				}
			}
		}
	} else {
		// 管理员随意设置
		// 根据menuIds查询查询菜单
		menus, err := s.menuc.ListMenu(ctx)
		if err != nil {
			return nil, errors.New("获取菜单列表失败" + err.Error())
		}
		for _, menuId := range menuIds {
			for _, menu := range menus {
				if menuId == menu.MenuId {
					reqMenus = append(reqMenus, menu)
				}
			}
		}
	}
	list[0].Menus = reqMenus
	err = s.rolec.UpdateRoleMenus(ctx, list[0])
	return &v1.CommonRes{}, err
}

//获取角色的权限接口
func (s *SysService) GetApisByRoleId(ctx context.Context, in *v1.QueryApisByRoleIdReq) (*v1.QueryApisByRoleIdRes, error) {
	if in.RoleId < 1 {
		return nil, errors.New("角色ID必须大于0")
	}
	reply, err := s.rolec.GetRoleApisList(ctx, []int64{in.RoleId})
	if err != nil {
		return nil, err
	}
	if len(reply) == 0 {
		return nil, errors.New("未获取到角色信息")
	}
	roleKey := reply[0].RoleKey
	apisList, err := s.rolec.GetRoleApisByRoleKeyword(roleKey)
	if err != nil {
		return nil, err
	}
	list := make([]*v1.QueryApisByRoleIdData, len(apisList))
	for k, v := range apisList {
		info := &v1.QueryApisByRoleIdData{
			ApiId:     v.ApiId,
			Handle:    v.Handle,
			Title:     v.Title,
			Path:      v.Path,
			Method:    v.Method,
			Category:  v.Category,
			CreatedAt: global.GetDateByUnix(v.CreatedAt),
			UpdatedAt: global.GetDateByUnix(v.UpdatedAt),
			CreateBy:  v.CreateBy,
			UpdateBy:  v.UpdateBy,
		}
		list[k] = info
	}
	return &v1.QueryApisByRoleIdRes{List: list, Total: int64(len(list))}, nil
}

//更新角色的权限接口
func (s *SysService) UpdateRoleApisById(ctx context.Context, in *v1.UpdateApisRoleReq) (*v1.CommonRes, error) {
	if in.RoleId < 1 {
		return nil, errors.New("角色ID必须大于0")
	}

	// 根据path中的角色ID获取该角色信息
	roleList, err := s.rolec.GetInfo(ctx, []int64{in.RoleId})
	if err != nil {
		return nil, err
	}
	if len(roleList) == 0 {
		return nil, errors.New("未获取到角色信息")
	}
	roleKey := roleList[0].RoleKey
	// 当前用户角色排序最小值（最高等级角色）以及当前用户
	ctxUser, err := s.userc.GetInfoById(ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	// (非管理员)不能更新比自己角色等级高或相等角色的权限接口

	// 获取当前用户所拥有的权限接口
	ctxRoles := ctxUser.Roles
	ctxRolesPolicies := make([][]string, 0)
	for _, role := range ctxRoles {
		policy := s.rolec.GetFilteredPolicy(role.RoleKey)
		ctxRolesPolicies = append(ctxRolesPolicies, policy...)
	}
	// 得到path中的角色ID对应角色能够设置的权限接口集合
	for _, policy := range ctxRolesPolicies {
		policy[0] = roleKey
	}
	// 生成前端想要设置的角色policies
	reqRolePolicies := make([][]string, 0)
	if len(in.ApiIds) > 0 {
		// 前端传来最新的ApiID集合
		apis, err := s.apic.GetApisById(ctx, in.ApiIds)
		if err != nil {
			return nil, err
		}
		for _, api := range apis {
			reqRolePolicies = append(reqRolePolicies, []string{
				roleKey, api.Path, api.Method,
			})
		}
	}

	// (非管理员)不能把角色的权限接口设置的比当前用户所拥有的权限接口多
	if in.UserId != 1 {
		for _, reqPolicy := range reqRolePolicies {
			if !funk.Contains(ctxRolesPolicies, reqPolicy) {
				return nil, errors.New(fmt.Sprintf("无权设置路径为%s,请求方式为%s的接口", reqPolicy[1], reqPolicy[2]))
			}
		}
	}
	// 更新角色的权限接口
	err = s.rolec.UpdateRoleApis(roleKey, reqRolePolicies)
	if err != nil {
		return nil, err
	}
	return &v1.CommonRes{}, nil
}
