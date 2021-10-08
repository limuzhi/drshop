package biz

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/global"
	"drpshop/internal/apps/sys/model"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/thoas/go-funk"
	"sort"
)

type SysMenuRepo interface {
	//列表
	ListMenu(ctx context.Context) ([]*model.SysMenu, error)
	//添加
	CreateMenu(ctx context.Context, in *model.SysMenu) error
	//修改
	UpdateMenu(ctx context.Context, id int64, in *model.SysMenu) error
	//获取
	GetMenuInfo(ctx context.Context, id int64) (*model.SysMenu, error)
	//批量删除
	BatchDeleteByIds(ctx context.Context, ids []int64) error
}

type SysMenuUsecase struct {
	repo     SysMenuRepo
	userRepo SysUserRepo
	roleRepo SysRoleRepo
	log      *log.Helper
}

func NewSysMenuUsecase(repo SysMenuRepo, userRepo SysUserRepo, roleRepo SysRoleRepo, logger log.Logger) *SysMenuUsecase {
	return &SysMenuUsecase{
		repo:     repo,
		userRepo: userRepo,
		roleRepo: roleRepo,
		log:      log.NewHelper(log.With(logger, "module", "sys/biz/sys_menu")),
	}
}

//列表
func (uc *SysMenuUsecase) ListMenu(ctx context.Context) ([]*model.SysMenu, error) {
	return uc.repo.ListMenu(ctx)
}

//添加
func (uc *SysMenuUsecase) CreateMenu(ctx context.Context, in *model.SysMenu) error {
	return uc.repo.CreateMenu(ctx, in)
}

//修改
func (uc *SysMenuUsecase) UpdateMenu(ctx context.Context, id int64, in *model.SysMenu) error {
	return uc.repo.UpdateMenu(ctx, id, in)
}

//获取
func (uc *SysMenuUsecase) GetMenuInfo(ctx context.Context, id int64) (*model.SysMenu, error) {
	return uc.repo.GetMenuInfo(ctx, id)
}

//批量删除
func (uc *SysMenuUsecase) BatchDeleteByIds(ctx context.Context, ids []int64) error {
	return uc.repo.BatchDeleteByIds(ctx, ids)
}

//TODO
//获取菜单树
func (uc *SysMenuUsecase) ListMenuTree(ctx context.Context) ([]*v1.MenuListData, error) {
	list, err := uc.repo.ListMenu(ctx)
	if err != nil {
		return nil, err
	}
	return uc.GenMenuTreeToMenuListData(0, list), nil
}

// 根据用户ID获取用户的权限(可访问)菜单列表
func (uc *SysMenuUsecase) GetUserMenusByUserId(ctx context.Context, uid int64) ([]*model.SysMenu, error) {
	userInfo, err := uc.userRepo.GetInfoById(ctx, uid)
	if err != nil {
		return nil, err
	}
	// 获取角色
	roleIds := make([]int64, 0)
	for _, role := range userInfo.Roles {
		roleIds = append(roleIds, role.RoleId)
	}

	allRoleList, err := uc.roleRepo.GetRoleMenuList(ctx, roleIds)
	if err != nil {
		return nil, err
	}
	// 所有角色的菜单集合
	allRoleMenus := make([]*model.SysMenu, 0)
	for _, userRole := range allRoleList {
		// 获取角色的菜单
		menus := userRole.Menus
		allRoleMenus = append(allRoleMenus, menus...)
	}
	// 所有角色的菜单集合去重
	allRoleMenusId := make([]int, 0)
	for _, menu := range allRoleMenus {
		allRoleMenusId = append(allRoleMenusId, int(menu.MenuId))
	}

	allRoleMenusIdUniq := funk.UniqInt(allRoleMenusId)
	allRoleMenusUniq := make([]*model.SysMenu, 0)
	for _, id := range allRoleMenusIdUniq {
		for _, menu := range allRoleMenus {
			if id == int(menu.MenuId) {
				allRoleMenusUniq = append(allRoleMenusUniq, menu)
				break
			}
		}
	}
	// 获取状态status为正常
	accessMenus := make([]*model.SysMenu, 0)
	for _, menu := range allRoleMenusUniq {
		if menu.Status == int(model.Enabled) {
			accessMenus = append(accessMenus, menu)
		}
	}
	sort.Sort(model.SysMenuSlice(accessMenus))
	return accessMenus, err
}

// 根据用户ID获取用户的权限(可访问)菜单树
func (uc *SysMenuUsecase) GetUserMenuTreeByUserId(ctx context.Context, uid int64) ([]*v1.MenuListData, error) {
	menuList, err := uc.GetUserMenusByUserId(ctx, uid)
	if err != nil {
		return nil, err
	}
	tree := uc.GenMenuTreeToMenuListData(0, menuList)
	return tree, err
}

func (uc *SysMenuUsecase) GenMenuTreeToMenuListData(parentId int64, menus []*model.SysMenu) []*v1.MenuListData {
	tree := make([]*v1.MenuListData, 0)
	for _, m := range menus {
		menuInfo := uc.DtoOut(m)
		if m.Pid == parentId {
			children := uc.GenMenuTreeToMenuListData(m.MenuId, menus)
			menuInfo.Children = children
			tree = append(tree, menuInfo)
		}
	}
	return tree
}

func (uc *SysMenuUsecase) GetMenusListDto(menus []*model.SysMenu) []*v1.MenuListData {
	out := make([]*v1.MenuListData, len(menus))
	for k, m := range menus {
		menuInfo := uc.DtoOut(m)
		out[k] = menuInfo
	}
	return out
}

func (uc *SysMenuUsecase) DtoOut(data *model.SysMenu) *v1.MenuListData {
	out := &v1.MenuListData{
		MenuId:     data.MenuId,
		Pid:        data.Pid,
		Name:       data.Name,
		Title:      data.Title,
		Icon:       data.Icon,
		Sort:       int64(data.Sort),
		Status:     int64(data.Status),
		Path:       data.Path,
		Hidden:     int64(data.Hidden),
		NoCache:    int64(data.NoCache),
		JumpPath:   data.JumpPath,
		Component:  data.Component,
		Redirect:   data.Redirect,
		ActiveMenu: data.ActiveMenu,
		IsFrame:    int64(data.IsFrame),
		ModuleType: data.ModuleType,
		ModelId:    int64(data.ModelId),
		AlwaysShow: int64(data.AlwaysShow),
		Breadcrumb: int64(data.Breadcrumb),
		CreatedAt:  global.GetDateByUnix(data.CreatedAt),
		UpdatedAt:  global.GetDateByUnix(data.UpdatedAt),
		CreateBy:   data.CreateBy,
		UpdateBy:   data.UpdateBy,
	}
	return out
}
