package data

import (
	"context"
	"errors"
	"fmt"
	"strings"

	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/biz"
	"drpshop/internal/apps/sys/data/model"

	"github.com/go-kratos/kratos/v2/log"
)

type sysRoleRepo struct {
	data *Data
	log  *log.Helper
}

func NewSysRoleRepo(data *Data, logger log.Logger) biz.SysRoleRepo {
	return &sysRoleRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "sys/data/sys_user")),
	}
}

//获取
func (r *sysRoleRepo) GetInfo(ctx context.Context, roleIds []int64) ([]*model.SysRole, error) {
	var list []*model.SysRole
	err := r.data.db.WithContext(ctx).Where("role_id IN (?)", roleIds).Find(&list).Error
	return list, err
}

func (r *sysRoleRepo) CreateRole(ctx context.Context, in *model.SysRole) error {
	return r.data.db.WithContext(ctx).Create(in).Error
}

//修改
func (r *sysRoleRepo) UpdateRole(ctx context.Context, in *model.SysRole) error {
	tx := r.data.db.WithContext(ctx)
	err := tx.Model(&model.SysRole{}).Where("role_id = ?", in.RoleId).Updates(in).Error
	if err != nil {
		return err
	}
	return nil
}

//修改状态
func (r *sysRoleRepo) ChangeRoleStatus(ctx context.Context, id, status int64) error {
	tx := r.data.db.WithContext(ctx)
	return tx.Model(&model.SysRole{}).Where("role_id = ?", id).
		Update("status", status).Error
}

//列表
func (r *sysRoleRepo) ListRole(ctx context.Context, in *v1.RoleListReq) ([]*model.SysRole, int64, error) {
	tx := r.data.db.WithContext(ctx).Model(&model.SysRole{})
	roleName := strings.TrimSpace(in.Name)
	if roleName != "" {
		tx = tx.Where("name LIKE ?", fmt.Sprintf("%%%s%%", roleName))
	}
	roleKey := strings.TrimSpace(in.RoleKey)
	if roleKey != "" {
		tx = tx.Where("role_key LIKE ?", fmt.Sprintf("%%%s%%", roleKey))
	}
	if in.Status != 0 {
		tx = tx.Where("status = ?", in.Status)
	}
	var list []*model.SysRole
	var count int64
	pageIndex := int(in.PageNum)
	pageSize := int(in.PageSize)
	offset := (pageIndex - 1) * pageSize
	err := tx.Order("sort DESC").Offset(offset).Limit(pageSize).Find(&list).
		Limit(-1).Offset(-1).Count(&count).Error
	return list, count, err
}

//批量删除用户数据
func (r *sysRoleRepo) BatchDeleteRoleByIds(ctx context.Context, ids []int64) error {
	list, err := r.GetInfo(ctx, ids)
	if err != nil {
		return err
	}
	err = r.data.db.WithContext(ctx).Model(&model.SysRole{}).
		Select("Users", "Depts", "Menus").Unscoped().Delete(&list).Error
	// 删除成功就删除casbin policy
	if err == nil {
		for _, role := range list {
			roleKeyword := role.RoleKey
			rmPolicies := r.data.casbinEnforcer.GetFilteredPolicy(0, roleKeyword)
			if len(rmPolicies) > 0 {
				isRemoved, _ := r.data.casbinEnforcer.RemovePolicies(rmPolicies)
				if !isRemoved {
					return errors.New("删除角色成功, 删除角色关联权限接口失败")
				}
			}
		}
	}
	return err
}

//获取角色对应的菜单
func (r *sysRoleRepo) GetRoleMenuList(ctx context.Context, roleIds []int64) ([]*model.SysRole, error) {
	if len(roleIds) == 0 {
		return nil, errors.New("无角色")
	}
	var list []*model.SysRole
	err := r.data.db.WithContext(ctx).Where("role_id IN (?)", roleIds).
		Preload("Menus").Find(&list).Error
	return list, err
}

// 根据角色关键字获取角色的权限接口
func (r *sysRoleRepo) GetRoleApisByRoleKeyword(roleKey string) ([]*model.SysApis, error) {
	policies := r.data.casbinEnforcer.GetFilteredPolicy(0, roleKey)
	// 获取所有接口
	var apis []*model.SysApis
	err := r.data.db.Model(&model.SysApis{}).Find(&apis).Error
	if err != nil {
		return nil, errors.New("获取角色的权限接口失败")
	}
	accessApis := make([]*model.SysApis, 0)
	for _, policy := range policies {
		path := policy[1]
		method := policy[2]
		for _, api := range apis {
			if path == api.Path && method == api.Method {
				accessApis = append(accessApis, api)
				break
			}
		}
	}
	return accessApis, err
}

//根据角色ID获取角色
func (r *sysRoleRepo) GetRoleApisList(ctx context.Context, roleIds []int64) ([]*model.SysRole, error) {
	var list []*model.SysRole
	err := r.data.db.Model(&model.SysRole{}).Where("role_id IN (?)", roleIds).Find(&list).Error
	return list, err
}

// 更新角色的权限接口（先全部删除再新增）
func (r *sysRoleRepo) UpdateRoleApis(roleKeyword string, reqRolePolicies [][]string) error {
	// 先获取path中的角色ID对应角色已有的police(需要先删除的)
	err := r.data.casbinEnforcer.LoadPolicy()
	if err != nil {
		return errors.New("角色的权限接口策略加载失败")
	}
	rmPolicies := r.data.casbinEnforcer.GetFilteredPolicy(0, roleKeyword)
	if len(rmPolicies) > 0 {
		isRemoved, _ := r.data.casbinEnforcer.RemovePolicies(rmPolicies)
		if !isRemoved {
			return errors.New("更新角色的权限接口失败")
		}
	}
	if len(reqRolePolicies) > 0 {
		isAdded, _ := r.data.casbinEnforcer.AddPolicies(reqRolePolicies)
		if !isAdded {
			return errors.New("更新角色的权限接口失败")
		}
	}
	err = r.data.casbinEnforcer.LoadPolicy()
	if err != nil {
		return errors.New("更新角色的权限接口成功，角色的权限接口策略加载失败")
	} else {
		return err
	}
}

//列表
func (r *sysRoleRepo) AllListRole(ctx context.Context) ([]*model.SysRole, error) {
	var list []*model.SysRole
	err := r.data.db.WithContext(ctx).Model(&model.SysRole{}).
		Where("status = ?", model.Enabled).Order("sort asc").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *sysRoleRepo) ListRoleByIds(ctx context.Context, ids []int64) ([]*model.SysRole, error) {
	var list []*model.SysRole
	err := r.data.db.WithContext(ctx).Model(&model.SysRole{}).Where("role_id IN(?)", ids).
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *sysRoleRepo) UpdateRoleMenus(ctx context.Context, in *model.SysRole) error {
	return r.data.db.WithContext(ctx).Model(in).Omit("Menus.*").Association("Menus").Replace(in.Menus)
}

func (r *sysRoleRepo) GetFilteredPolicy(roleKey string) [][]string {
	return r.data.casbinEnforcer.GetFilteredPolicy(0, roleKey)
}

func (r *sysRoleRepo) CheckCasbinRole(subs []string, obj string, act string) bool {
	isPass := false
	for _, sub := range subs {
		pass, _ := r.data.casbinEnforcer.Enforce(sub, obj, act)
		if pass {
			isPass = true
			break
		}
	}
	return isPass
}

// 根据角色关键字获取角色的权限接口
func (r *sysRoleRepo) GetPermissionsByRoleKey(roleKey string) ([]string, error) {
	policies := r.data.casbinEnforcer.GetFilteredPolicy(0, roleKey)
	// 获取所有接口
	var apis []*model.SysApis
	err := r.data.db.Model(&model.SysApis{}).Find(&apis).Error
	if err != nil {
		return nil, errors.New("获取角色的权限接口失败")
	}
	permissionsList := make([]string, 0)
	for _, policy := range policies {
		path := policy[1]
		method := policy[2]
		for _, api := range apis {
			if path == api.Path && method == api.Method {
				if api.Permission != "" {
					permissionsList = append(permissionsList, api.Permission)
				}
				break
			}
		}
	}
	return permissionsList, err
}
