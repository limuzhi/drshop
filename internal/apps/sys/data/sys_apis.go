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

type sysApisRepo struct {
	data *Data
	log  *log.Helper
}

func NewSysApisRepo(data *Data, logger log.Logger) biz.SysApisRepo {
	return &sysApisRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "sys/data/sys_apis")),
	}
}

//列表
func (r *sysApisRepo) ListApis(ctx context.Context, in *v1.ApisListReq) ([]*model.SysApis, int64, error) {
	table := r.data.db.WithContext(ctx)
	table = table.Model(&model.SysApis{})
	method := strings.TrimSpace(in.Method)
	if method != "" {
		table = table.Where("method = ?", method)
	}
	path := strings.TrimSpace(in.Path)
	if path != "" {
		table = table.Where("path LIKE ?", fmt.Sprintf("%%%s%%", path))
	}
	category := strings.TrimSpace(in.Category)
	if category != "" {
		table = table.Where("category LIKE ?", fmt.Sprintf("%%%s%%", category))
	}
	var list []*model.SysApis
	var total int64
	err := table.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	pageNum := int(in.PageInfo.PageNum)
	pageSize := int(in.PageInfo.PageSize)
	err = table.Offset((pageNum - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&list).Error
	return list, total, err
}

//创建
func (r *sysApisRepo) CreateApis(ctx context.Context, in *model.SysApis) error {
	return r.data.db.WithContext(ctx).Create(in).Error
}

//更新
func (r *sysApisRepo) UpdateApis(ctx context.Context, apiId int64, in *model.SysApis) error {
	// 根据id获取接口信息
	var oldApi model.SysApis
	tx := r.data.db.WithContext(ctx)
	err := tx.First(&oldApi, apiId).Error
	if err != nil {
		return errors.New("根据接口ID获取接口信息失败")
	}
	err = tx.Model(in).Where("api_id = ?", apiId).Updates(in).Error
	if err != nil {
		return err
	}
	// 更新了method和path就更新casbin中policy
	if oldApi.Path != in.Path || oldApi.Method != in.Method {
		policies := r.data.casbinEnforcer.GetFilteredPolicy(1, oldApi.Path, oldApi.Method)
		// 接口在casbin的policy中存在才进行操作
		if len(policies) > 0 {
			// 先删除
			isRemoved, _ := r.data.casbinEnforcer.RemovePolicies(policies)
			if !isRemoved {
				return errors.New("更新权限接口失败")
			}
			for _, policy := range policies {
				policy[1] = in.Path
				policy[2] = in.Method
			}
			// 新增
			isAdded, _ := r.data.casbinEnforcer.AddPolicies(policies)
			if !isAdded {
				return errors.New("更新权限接口失败")
			}
			// 加载policy
			err := r.data.casbinEnforcer.LoadPolicy()
			if err != nil {
				return errors.New("更新权限接口成功，权限接口策略加载失败")
			} else {
				return err
			}
		}
	}
	return err
}

//根据接口ID获取接口列表
func (r *sysApisRepo) GetApisById(ctx context.Context, ids []int64) ([]*model.SysApis, error) {
	var apis []*model.SysApis
	err := r.data.db.WithContext(ctx).Where("api_id IN (?)", ids).Find(&apis).Error
	if err != nil {
		return nil, err
	}
	return apis, nil
}

//批量删除接口
func (r *sysApisRepo) BatchDeleteByIds(ctx context.Context, ids []int64) error {
	apis, err := r.GetApisById(ctx, ids)
	if err != nil {
		return errors.New("根据接口ID获取接口列表失败")
	}
	if len(apis) == 0 {
		return errors.New("根据接口ID未获取到接口列表")
	}

	err = r.data.db.WithContext(ctx).Where("api_id IN (?)", ids).Unscoped().Delete(&model.SysApis{}).Error
	// 如果删除成功，删除casbin中policy
	if err == nil {
		for _, api := range apis {
			policies := r.data.casbinEnforcer.GetFilteredPolicy(1, api.Path, api.Method)
			if len(policies) > 0 {
				isRemoved, _ := r.data.casbinEnforcer.RemovePolicies(policies)
				if !isRemoved {
					return errors.New("删除权限接口失败")
				}
			}
		}
		// 重新加载策略
		err := r.data.casbinEnforcer.LoadPolicy()
		if err != nil {
			return errors.New("删除权限接口成功，权限接口策略加载失败")
		} else {
			return err
		}
	}
	return err
}

//根据接口路径和请求方式获取接口描述
func (r *sysApisRepo) GetApiDescByPath(ctx context.Context, path string, method string) (string, error) {
	var info model.SysApis
	err := r.data.db.WithContext(ctx).Where("path = ? AND method = ?", path, method).First(&info).Error
	if err != nil {
		return "", err
	}
	return info.Title, nil
}

//获取接口树(按接口Category字段分类)
func (r *sysApisRepo) GetListOrderCategory(ctx context.Context) ([]*model.SysApis, error) {
	var list []*model.SysApis
	err := r.data.db.WithContext(ctx).Order("category").
		Order("created_at").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

//查找
func (r *sysApisRepo) GetListByPath(ctx context.Context, path []string) ([]*model.SysApis, error) {
	var list []*model.SysApis
	err := r.data.db.WithContext(ctx).Model(&model.SysApis{}).
		Select("path,method,title").Where("path IN(?) ", path).First(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}
