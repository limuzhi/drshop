package data

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/biz"
	"drpshop/internal/apps/sys/model"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"strings"
)

type sysConfigRepo struct {
	data *Data
	log  *log.Helper
}

func NewSysConfigRepo(data *Data, logger log.Logger) biz.SysConfigRepo {
	return &sysConfigRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "sys/data/sys_config")),
	}
}

//列表
func (r *sysConfigRepo) ListConfig(ctx context.Context, in *v1.ConfigListReq) ([]*model.SysConfig, int64, error) {
	table := r.data.db.WithContext(ctx)
	table = table.Model(&model.SysConfig{})
	if in != nil {
		configName := strings.TrimSpace(in.ConfigName)
		if configName != "" {
			table = table.Where("config_name LIKE ?", "%"+configName+"%")
		}
		if in.ConfigType != "" {
			configType, _ := strconv.Atoi(in.ConfigType)
			table = table.Where("config_type = ?", configType)
		}
		configKey := strings.TrimSpace(in.ConfigKey)
		if configKey != "" {
			table = table.Where("config_key like ?", "%"+configKey+"%")
		}
	}

	var list []*model.SysConfig
	var total int64
	err := table.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	pageNum := int(in.PageInfo.PageNum)
	pageSize := int(in.PageInfo.PageSize)
	err = table.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

//删除
func (r *sysConfigRepo) BatchDeleteByIds(ctx context.Context, ids []int64) error {
	err := r.data.db.WithContext(ctx).Where("config_id IN (?)", ids).Unscoped().Delete(&model.SysConfig{}).Error
	return err
}

//创建
func (r *sysConfigRepo) CreateConfig(ctx context.Context, in *model.SysConfig) error {
	return r.data.db.WithContext(ctx).Create(in).Error
}

//修改
func (r *sysConfigRepo) UpdateConfig(ctx context.Context, id int64, in *model.SysConfig) error {
	var oldInfo model.SysConfig
	tx := r.data.db.WithContext(ctx)
	err := tx.Where("config_id = ?", id).First(&oldInfo).Error
	if err != nil {
		return errors.New("根据ID获取信息失败")
	}
	err = tx.Where("config_id = ?", id).Updates(in).Error
	if err != nil {
		return err
	}
	return err
}

//获取
func (r *sysConfigRepo) GetInfoById(ctx context.Context, id int64) (*model.SysConfig, error) {
	var info *model.SysConfig
	err := r.data.db.WithContext(ctx).Where("config_id =", id).First(&info).Error
	if err != nil {
		return nil, err
	}
	return info, nil
}

//获取
func (r *sysConfigRepo) GetInfoByConfigKey(ctx context.Context, configKey string) (*model.SysConfig, error) {
	var info *model.SysConfig
	err := r.data.db.WithContext(ctx).Where("config_key = ?", configKey).First(&info).Error
	if err != nil {
		return nil, err
	}
	return info, nil
}
