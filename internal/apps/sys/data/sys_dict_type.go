package data

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/biz"
	"drpshop/internal/apps/sys/data/model"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"strings"
)

type sysDictTypeRepo struct {
	data *Data
	log  *log.Helper
}

func NewSysDictTypeRepo(data *Data, logger log.Logger) biz.SysDictTypeRepo {
	return &sysDictTypeRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "sys/data/sys_apis")),
	}
}

//列表
func (r *sysDictTypeRepo) ListDictType(ctx context.Context, in *v1.DictTypeListReq) ([]*model.SysDictType, int64, error) {
	table := r.data.db.WithContext(ctx)
	table = table.Model(&model.SysDictType{})
	if in != nil {
		dictName := strings.TrimSpace(in.DictName)
		if dictName != "" {
			table = table.Where("dict_name LIKE ?", fmt.Sprintf("%%%s%%", dictName))
		}
		dictType := strings.TrimSpace(in.DictType)
		if dictType != "" {
			table = table.Where("dict_type LIKE ?", fmt.Sprintf("%%%s%%", dictType))
		}
		if in.Status != 0 {
			table = table.Where("status = ?", in.Status)
		}
	}

	var list []*model.SysDictType
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
func (r *sysDictTypeRepo) CreateDictType(ctx context.Context, in *model.SysDictType) error {
	return r.data.db.WithContext(ctx).Create(in).Error
}

//修改
func (r *sysDictTypeRepo) UpdateDictType(ctx context.Context, id int64, in *model.SysDictType) error {
	//事务
	return r.data.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var info *model.SysDictType
		if err := tx.Where("dict_id = ?", id).First(&info).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.SysDictType{}).
			Where("dict_id = ?", id).Updates(in).Error; err != nil {
			return err
		}
		if in.DictType != info.DictType {
			if err := tx.Where("dict_type = ?", info.DictType).
				Update("dict_type", in.DictType).Error; err != nil {
				return err
			}
		}
		// 返回 nil 提交事务
		return nil
	})
}

//获取
func (r *sysDictTypeRepo) GetInfoById(ctx context.Context, id int64) (*model.SysDictType, error) {
	var info *model.SysDictType
	err := r.data.db.WithContext(ctx).Where("dict_id = ?", id).First(&info).Error
	return info, err
}

//批量删除接口
func (r *sysDictTypeRepo) BatchDeleteByIds(ctx context.Context, ids []int64) error {
	var err error
	var list []*model.SysDictType
	if err = r.data.db.WithContext(ctx).Where("dict_id IN(?)", ids).Find(&list).Error; err != nil {
		return err
	}
	dictTypeArr := make([]string, 0)
	for _, v := range list {
		dictTypeArr = append(dictTypeArr, v.DictType)
	}
	if len(dictTypeArr) > 0 {
		return r.data.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("dict_id IN(?)", ids).Unscoped().Delete(&model.SysDictType{}).Error; err != nil {
				return err
			}
			if err := tx.Where("dict_type IN(?)", dictTypeArr).Unscoped().Delete(&model.SysDictData{}).Error; err != nil {
				return err
			}
			return nil
		})
	}
	return err
}

//获取所有正常状态下的字典类型
func (r *sysDictTypeRepo) GetAllDictType(ctx context.Context) ([]*model.SysDictType, error) {
	var list []*model.SysDictType
	err := r.data.db.WithContext(ctx).Model(&model.SysDictType{}).
		Where("status = ?", model.Enabled).Find(&list).Error
	return list, err
}

//检查类型是否已经存在
func (r *sysDictTypeRepo) CheckByDictType(ctx context.Context, dictType string, ids ...int64) bool {
	table := r.data.db.WithContext(ctx).Model(&model.SysDictType{}).Where("dict_type = ?", dictType)
	if len(ids) > 0 {
		table = table.Where("dict_id = ?", ids)
	}
	var info *model.SysDictType
	if err := table.First(&info).Error; err != nil {
		return false
	}
	return true
}
