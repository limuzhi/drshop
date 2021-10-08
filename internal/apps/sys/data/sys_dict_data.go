package data

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/biz"
	"drpshop/internal/apps/sys/model"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"strings"
)

type sysDictDataRepo struct {
	data *Data
	log  *log.Helper
}

func NewSysDictDataRepo(data *Data, logger log.Logger) biz.SysDictDataRepo {
	return &sysDictDataRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "sys/data/sys_dict_data")),
	}
}

//列表
func (r *sysDictDataRepo) ListDictData(ctx context.Context, in *v1.DictDataListReq) ([]*model.SysDictData, int64, error) {
	table := r.data.db.WithContext(ctx)
	table = table.Model(&model.SysDictData{})
	if in != nil {
		dictLabel := strings.TrimSpace(in.DictLabel)
		if dictLabel != "" {
			table = table.Where("dict_label LIKE ?", fmt.Sprintf("%%%s%%", dictLabel))
		}
		dictType := strings.TrimSpace(in.DictType)
		if dictType != "" {
			table = table.Where("dict_type = ?", dictType)
		}
		if in.Status != 0 {
			table = table.Where("status = ?", in.Status)
		}
	}
	var list []*model.SysDictData
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

//添加
func (r *sysDictDataRepo) CreateDictData(ctx context.Context, in *model.SysDictData) error {
	return r.data.db.WithContext(ctx).Create(in).Error
}

//修改
func (r *sysDictDataRepo) UpdateDictData(ctx context.Context, id int64, in *model.SysDictData) error {
	err := r.data.db.WithContext(ctx).Model(&model.SysDictData{}).
		Where("dict_code = ?", id).Updates(in).Error
	return err
}

//获取
func (r *sysDictDataRepo) GetInfoById(ctx context.Context, id int64) (*model.SysDictData, error) {
	var info *model.SysDictData
	err := r.data.db.WithContext(ctx).Where("dict_code = ?", id).First(&info).Error
	return info, err
}

//批量删除接口
func (r *sysDictDataRepo) BatchDeleteByIds(ctx context.Context, ids []int64) error {
	err := r.data.db.WithContext(ctx).
		Where("dict_code IN (?)", ids).Unscoped().Delete(&model.SysDictData{}).Error
	return err
}

//获取列表 dictType
func (r *sysDictDataRepo) ListByDictType(ctx context.Context, dictType string) ([]*model.SysDictData, error) {
	var list []*model.SysDictData
	err := r.data.db.WithContext(ctx).
		Where("dict_type = ?", dictType).Order("dict_sort asc").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, err
}
