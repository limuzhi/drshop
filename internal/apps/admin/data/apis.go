package data

import (
	"context"
	"drpshop/api/common"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/admin/biz"
	"drpshop/internal/apps/admin/vo"
	"github.com/go-kratos/kratos/v2/log"
	"strings"
)

type apisRepo struct {
	data *Data
	log  *log.Helper
}

func NewApisRepo(data *Data, logger log.Logger) biz.ApisRepo {
	return &apisRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "admin/data/apis")),
	}
}

func (r *apisRepo) GetApiTreeList(ctx context.Context) (*v1.ApisTreeListRes, error) {
	return r.data.sc.ApisTreeList(ctx, &common.NullReq{})
}

func (r *apisRepo) GetApiList(ctx context.Context, in *vo.ApiListReq) (*v1.ApisListRes, error) {
	return r.data.sc.ApisList(ctx, &v1.ApisListReq{
		Method:   strings.TrimSpace(in.Method),
		Path:     strings.TrimSpace(in.Path),
		Category: strings.TrimSpace(in.Category),
		PageInfo: &common.PageReq{
			PageNum:  int64(in.GetPageIndex()),
			PageSize: int64(in.GetPageSize()),
		},
	})
}
func (r *apisRepo) ApiCreate(ctx context.Context, in *vo.CreateApiReq) error {
	_, err := r.data.sc.ApisAdd(ctx, &v1.ApisAddReq{
		Handle:     strings.TrimSpace(in.Handle),
		Title:      strings.TrimSpace(in.Title),
		Path:       strings.TrimSpace(in.Path),
		Method:     strings.TrimSpace(in.Method),
		Category:   strings.TrimSpace(in.Category),
		Permission: strings.TrimSpace(in.Permission),
		CreateBy:   in.CreateBy,
	})
	return err
}

func (r *apisRepo) ApiUpdate(ctx context.Context, in *vo.UpdateApiReq) error {
	_, err := r.data.sc.ApisUpdate(ctx, &v1.ApisUpdateReq{
		ApiId:      in.ApiId,
		Handle:     strings.TrimSpace(in.Handle),
		Title:      strings.TrimSpace(in.Title),
		Path:       strings.TrimSpace(in.Path),
		Method:     strings.TrimSpace(in.Method),
		Category:   strings.TrimSpace(in.Category),
		Permission: strings.TrimSpace(in.Permission),
		UpdateBy:   in.UpdateBy,
	})
	return err
}

func (r *apisRepo) ApiBatchDelete(ctx context.Context, ids []int64) error {
	_, err := r.data.sc.ApisDelete(ctx, &v1.ApisDeleteReq{ApiIds: ids})
	return err
}
