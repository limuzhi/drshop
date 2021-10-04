package data

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/biz"
	"drpshop/internal/apps/sys/data/model"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"strings"
)

type sysPostRepo struct {
	data *Data
	log  *log.Helper
}

func NewSysPostRepo(data *Data, logger log.Logger) biz.SysPostRepo {
	return &sysPostRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "sys/data/sys_apis")),
	}
}

//列表
func (r *sysPostRepo) ListPost(ctx context.Context, in *v1.PostListReq) ([]*model.SysPost, int64, error) {
	table := r.data.db.WithContext(ctx)
	table = table.Model(&model.SysPost{})
	if in != nil {
		postName := strings.TrimSpace(in.PostName)
		if postName != "" {
			table = table.Where("post_name LIKE ?", fmt.Sprintf("%%%s%%", postName))
		}
		postCode := strings.TrimSpace(in.PostCode)
		if postCode != "" {
			table = table.Where("post_code LIKE ?", fmt.Sprintf("%%%s%%", postCode))
		}
		if in.Status != 0 {
			table = table.Where("status = ?", in.Status)
		}
	}
	var list []*model.SysPost
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
func (r *sysPostRepo) CreatePost(ctx context.Context, in *model.SysPost) error {
	return r.data.db.WithContext(ctx).Create(in).Error
}

//修改
func (r *sysPostRepo) UpdatePost(ctx context.Context, id int64, in *model.SysPost) error {
	err := r.data.db.WithContext(ctx).Model(&model.SysPost{}).
		Where("post_id = ?", id).Updates(in).Error
	return err
}

//获取
func (r *sysPostRepo) GetInfoById(ctx context.Context, id int64) (*model.SysPost, error) {
	var info *model.SysPost
	err := r.data.db.WithContext(ctx).Where("post_id = ?", id).First(&info).Error
	return info, err
}

//删除
func (r *sysPostRepo) BatchDeleteByIds(ctx context.Context, ids []int64) error {
	var posts []*model.SysPost
	tx := r.data.db.WithContext(ctx)
	err := tx.Where("post_id IN (?)", ids).Find(&posts).Error
	if err != nil {
		return err
	}
	return tx.Select("Users").Unscoped().Delete(&posts).Error
}

func (r *sysPostRepo) AllListPost(ctx context.Context) ([]*model.SysPost, error) {
	var list []*model.SysPost
	err := r.data.db.WithContext(ctx).Model(&model.SysPost{}).
		Where("status = ?", model.Enabled).Order("post_sort asc").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *sysPostRepo) ListPostByIds(ctx context.Context, ids []int64) ([]*model.SysPost, error) {
	var List []*model.SysPost
	err := r.data.db.WithContext(ctx).Model(&model.SysPost{}).
		Where("post_id IN (?)", ids).Find(&List).Error
	return List, err
}
