package data

import (
	"context"
	"drpshop/internal/apps/sys/biz"
	"drpshop/internal/apps/sys/data/model"
	"github.com/go-kratos/kratos/v2/log"
)

type sysMenuRepo struct {
	data *Data
	log  *log.Helper
}

func NewSysMenuRepo(data *Data, logger log.Logger) biz.SysMenuRepo {
	return &sysMenuRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "sys/data/sys_menu")),
	}
}

func (r *sysMenuRepo) ListMenu(ctx context.Context) ([]*model.SysMenu, error) {
	var list []*model.SysMenu
	err := r.data.db.WithContext(ctx).Order("sort asc").Preload("RoleList").Find(&list).Error
	return list, err
}

//添加
func (r *sysMenuRepo) CreateMenu(ctx context.Context, in *model.SysMenu) error {
	return r.data.db.WithContext(ctx).Create(in).Error
}

//修改
func (r *sysMenuRepo) UpdateMenu(ctx context.Context, id int64, in *model.SysMenu) error {
	return r.data.db.WithContext(ctx).Model(&model.SysMenu{}).
		Where("menu_id = ?", id).Updates(in).Error
}

//获取
func (r *sysMenuRepo) GetMenuInfo(ctx context.Context, id int64) (*model.SysMenu, error) {
	var info *model.SysMenu
	err := r.data.db.WithContext(ctx).Where("menu_id = ?", id).
		Preload("ApiList").First(&info).Error
	return info, err
}

//批量删除
func (r *sysMenuRepo) BatchDeleteByIds(ctx context.Context, ids []int64) error {
	var menus []*model.SysMenu
	tx := r.data.db.WithContext(ctx)
	err := tx.Where("menu_id IN (?)", ids).Find(&menus).Error
	if err != nil {
		return err
	}
	err = tx.Select("RoleList").Unscoped().Delete(&menus).Error
	return err
}
