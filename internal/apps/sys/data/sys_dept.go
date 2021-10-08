package data

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/biz"
	"drpshop/internal/apps/sys/model"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/util/gconv"
	"strings"
)

type sysDeptRepo struct {
	data *Data
	log  *log.Helper
}

func NewSysDeptRepo(data *Data, logger log.Logger) biz.SysDeptRepo {
	return &sysDeptRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "sys/data/sys_dept")),
	}
}

//列表
func (r *sysDeptRepo) ListDept(ctx context.Context, in *v1.DeptListReq) ([]*model.SysDept, int64, error) {
	table := r.data.db.WithContext(ctx)
	table = table.Model(&model.SysDept{})
	if in != nil {
		deptName := strings.TrimSpace(in.DeptName)
		if deptName != "" {
			table = table.Where("dept_name LIKE ?", "%"+deptName+"%")
		}

		if in.Status != "" {
			table = table.Where("status = ?", gconv.Int(in.Status))
		}
	}
	var list []*model.SysDept
	var total int64
	err := table.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = table.Find(&list).Error
	return list, total, err
}

//删除
func (r *sysDeptRepo) BatchDeleteByIds(ctx context.Context, ids []int64) error {
	var deptList []*model.SysDept
	tx := r.data.db.WithContext(ctx)
	err := tx.Model(&model.SysDept{}).Where("dept_id IN (?)", ids).Find(&deptList).Error
	if err != nil {
		return err
	}
	return tx.Select("RoleList").Unscoped().Delete(&deptList).Error
}

//创建
func (r *sysDeptRepo) CreateDept(ctx context.Context, in *model.SysDept) error {
	return r.data.db.WithContext(ctx).Create(in).Error
}

//修改
func (r *sysDeptRepo) UpdateDept(ctx context.Context, id int64, in *model.SysDept) error {
	err := r.data.db.WithContext(ctx).Model(&model.SysDept{}).
		Where("dept_id = ?", id).Updates(in).Error
	return err
}

//获取
func (r *sysDeptRepo) GetInfoById(ctx context.Context, id int64) (*model.SysDept, error) {
	var info *model.SysDept
	err := r.data.db.WithContext(ctx).Where("dept_id = ?", id).First(&info).Error
	return info, err
}

//
func (r *sysDeptRepo) GetRoleDepts(ctx context.Context, roleId int64) ([]int64, error) {
	var list []*model.SysRoleDept
	err := r.data.db.WithContext(ctx).Model(&model.SysRoleDept{}).
		Where("role_id = ?", roleId).Find(&list).Error
	if err != nil {
		return nil, err
	}
	result := make([]int64, 0)
	for _, v := range list {
		result = append(result, v.DeptId)
	}
	return result, nil
}

func (r *sysDeptRepo) AllListDept(ctx context.Context) ([]*model.SysDept, error) {
	var list []*model.SysDept
	err := r.data.db.WithContext(ctx).Model(&model.SysDept{}).
		Where("status = ?", model.Enabled).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}
