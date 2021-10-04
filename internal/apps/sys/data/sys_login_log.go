package data

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/biz"
	"drpshop/internal/apps/sys/data/model"
	"github.com/go-kratos/kratos/v2/log"
)

type sysLoginLogRepo struct {
	data *Data
	log  *log.Helper
}

func NewSysLoginLogRepo(data *Data, logger log.Logger) biz.SysLoginLogRepo {
	return &sysLoginLogRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "sys/data/sys_login_log")),
	}
}

//列表
func (r *sysLoginLogRepo) ListLoginLog(ctx context.Context,in *v1.LoginlogListReq) ([]*model.SysLoginLog, int64, error) {
	tx := r.data.db.WithContext(ctx)
	tx = tx.Model(&model.SysLoginLog{})
	if in.LoginName != "" {
		tx = tx.Where("login_name like ?", "%"+in.LoginName+"%")
	}
	if in.Status != 0 {
		tx = tx.Where("status = ?", in.Status)
	}
	if in.Ipaddr != "" {
		tx = tx.Where("ipaddr like ?", "%"+in.Ipaddr+"%")
	}
	if in.LoginLocation != "" {
		tx = tx.Where("login_location like ?", "%"+in.LoginLocation+"%")
	}
	if in.PageInfo.BeginTime > 0 {
		tx = tx.Where("login_time >= ?", in.PageInfo.BeginTime)
	}
	if in.PageInfo.EndTime > 0 {
		tx = tx.Where("login_time <= ?", in.PageInfo.EndTime)
	}
	order := "login_time DESC"
	var list []*model.SysLoginLog
	var total int64
	err := tx.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	pageNum := int(in.PageInfo.PageNum)
	pageSize := int(in.PageInfo.PageSize)
	err = tx.Offset((pageNum - 1) * pageSize).Limit(pageSize).Order(order).Find(&list).Error
	return list, total, err
}

//创建
func (r *sysLoginLogRepo) CreateLoginLog(ctx context.Context, ins []*model.SysLoginLog) error {
	return r.data.db.WithContext(ctx).Create(ins).Error
}

//删除
func (r *sysLoginLogRepo) BatchDeleteByIds(ctx context.Context, ins []int64) error {
	return r.data.db.WithContext(ctx).Model(&model.SysLoginLog{}).
		Where("login_id IN(?)", ins).
		Delete(&model.SysLoginLog{}).Error
}

func (r *sysLoginLogRepo) GetInfo(ctx context.Context, id int64) (*model.SysLoginLog, error) {
	var info *model.SysLoginLog
	err := r.data.db.WithContext(ctx).Where("login_id = ?", id).First(&info).Error
	if err != nil {
		return nil, err
	}
	return info, nil
}


func (r *sysLoginLogRepo) ClearLoginLog(ctx context.Context) error {
	modelLog := &model.SysLoginLog{}
	return r.data.db.WithContext(ctx).Exec("truncate " + modelLog.TableName()).Error
}
