package data

import (
	"context"
	"drpshop/api/common"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/admin/biz"
	"drpshop/internal/apps/admin/vo"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/util/gconv"
	"strings"
)

type postRepo struct {
	data *Data
	log  *log.Helper
}

func NewPostRepo(data *Data, logger log.Logger) biz.PostRepo {
	return &postRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "admin/data/post")),
	}
}

//岗位
func (r *postRepo) PostDetail(ctx context.Context, id int64) (*v1.PostDetailRes, error) {
	return r.data.sc.PostDetail(ctx, &v1.PostDetailReq{PostId: id})
}

func (r *postRepo) PostList(ctx context.Context, in *vo.PostListReq) (*v1.PostListRes, error) {
	return r.data.sc.PostList(ctx, &v1.PostListReq{
		Status:   gconv.Int64(in.Status),
		PostName: strings.TrimSpace(in.PostName),
		PostCode: strings.TrimSpace(in.PostCode),
		PageInfo: &common.PageReq{
			PageNum: int64(in.GetPageIndex()),
			PageSize: int64(in.GetPageSize()),
		},
	})
}

func (r *postRepo) PostCreate(ctx context.Context, in *vo.PostCreateReq) error {
	data := &v1.PostCreateUpdateReq{
		PostCode: strings.TrimSpace(in.PostCode),
		PostName: strings.TrimSpace(in.PostName),
		PostSort: gconv.Int64(in.PostSort),
		Remark:   in.Remark,
		Status:   gconv.Int64(in.Status),
		CreateBy: in.CreateBy,
	}
	_, err := r.data.sc.PostCreate(ctx, data)
	return err
}

func (r *postRepo) PostUpdate(ctx context.Context, in *vo.PostUpdateReq) error {
	data := &v1.PostCreateUpdateReq{
		PostId:   in.PostId,
		PostCode: strings.TrimSpace(in.PostCode),
		PostName: strings.TrimSpace(in.PostName),
		PostSort: gconv.Int64(in.PostSort),
		Remark:   in.Remark,
		Status:   gconv.Int64(in.Status),
		UpdateBy: in.UpdateBy,
	}
	_, err := r.data.sc.PostUpdate(ctx, data)
	return err
}

func (r *postRepo) PostDelete(ctx context.Context, ids []int64) error {
	_, err := r.data.sc.PostDelete(ctx, &v1.PostDeleteReq{PostIds: ids})
	return err
}
