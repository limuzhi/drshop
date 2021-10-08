package data

import (
	"context"
	"drpshop/internal/apps/sys/biz"
	"drpshop/internal/apps/sys/model"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"strings"
)

type sysUserOnlineRepo struct {
	data *Data
	log  *log.Helper
}

func NewSysUserOnlineRepo(data *Data, logger log.Logger) biz.SysUserOnlineRepo {
	return &sysUserOnlineRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "sys/data/sys_user_online")),
	}
}

//列表
func (r *sysUserOnlineRepo) ListUserOnline(ctx context.Context, in *model.SysUserOnlineSearchReq) ([]*model.SysUserOnline, int64, error) {
	table := r.data.db.WithContext(ctx)
	table = table.Model(&model.SysUserOnline{})
	if in != nil {
		username := strings.TrimSpace(in.Username)
		if username != "" {
			table = table.Where("method LIKE ?", fmt.Sprintf("%%%s%%", username))
		}
		ip := strings.TrimSpace(in.Ip)
		if ip != "" {
			table = table.Where("ip LIKE ?", fmt.Sprintf("%%%s%%", ip))
		}
	}
	var list []*model.SysUserOnline
	var total int64
	err := table.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	if in.PageNum <= 0 {
		in.PageNum = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 10
	}
	pageNum := int(in.PageNum)
	pageSize := int(in.PageSize)
	err = table.Offset((pageNum - 1) * pageSize).Limit(pageSize).Order("create_time DESC").Find(&list).Error
	return list, total, err
}

//创建
func (r *sysUserOnlineRepo) CreateUserOnline(ctx context.Context, in *model.SysUserOnline) error {
	return r.data.db.WithContext(ctx).Create(in).Error
}

//根据ID获取
func (r *sysUserOnlineRepo) GetListByIds(ctx context.Context, ids []int64) ([]*model.SysUserOnline, error) {
	var list []*model.SysUserOnline
	err := r.data.db.WithContext(ctx).Where("id in(?)", ids).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

//删除用户在线状态操作
func (r *sysUserOnlineRepo) DeleteByToken(ctx context.Context, token string) error {
	err := r.data.db.WithContext(ctx).Where("token = ?", token).
		Unscoped().Delete(&model.SysUserOnline{}).Error
	return err
}

//删除用户在线状态操作
func (r *sysUserOnlineRepo) BatchDeleteByIds(ctx context.Context, ids []int64) ([]string, error) {
	list, err := r.GetListByIds(ctx, ids)
	if err != nil {
		return nil, err
	}
	tokens := make([]string, 0)
	for _, v := range list {
		tokens = append(tokens, v.Token)
	}
	err = r.data.db.WithContext(ctx).
		Where("id IN(?)", ids).Unscoped().Delete(&model.SysUserOnline{}).Error

	return tokens, err
}

//根据Token获取
func (r *sysUserOnlineRepo) SaveUserOnline(ctx context.Context, data *model.SysUserOnline) error {
	var info *model.SysUserOnline
	err := r.data.db.WithContext(ctx).Where("token = ? ", data.Token).First(&info).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			results := r.data.db.WithContext(ctx).Model(&model.SysUserOnline{}).Create(&data)
			if results.Error != nil {
				return results.Error
			}
		}
	} else {
		data.ID = info.ID
		results := r.data.db.WithContext(ctx).Model(info).
			Updates(data)
		if results.Error != nil {
			return results.Error
		}
	}
	return nil
}
