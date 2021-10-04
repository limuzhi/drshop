package data

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/admin/biz"
	"drpshop/internal/apps/admin/vo"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/util/gconv"
	"strings"
)

type deptRepo struct {
	data *Data
	log  *log.Helper
}

func NewDeptRepo(data *Data, logger log.Logger) biz.DeptRepo {
	return &deptRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "admin/data/dept")),
	}
}

func (r *deptRepo) DeptList(ctx context.Context, in *vo.DeptListReq) (*v1.DeptListRes, error) {
	return r.data.sc.DeptList(ctx, &v1.DeptListReq{
		Status:   strings.TrimSpace(in.Status),
		DeptName: strings.TrimSpace(in.DeptName),
	})
}

func (r *deptRepo) DeptTree(ctx context.Context, in *vo.DeptListReq) (*v1.DeptListRes, error) {
	return r.data.sc.DeptTree(ctx, &v1.DeptListReq{
		Status:   strings.TrimSpace(in.Status),
		DeptName: strings.TrimSpace(in.DeptName),
	})
}

func (r *deptRepo) DeptDetail(ctx context.Context, deptId int64) (*v1.DeptDetailRes, error) {
	return r.data.sc.DeptDetail(ctx, &v1.DeptDetailReq{DeptId: deptId})

}

func (r *deptRepo) DeptCreate(ctx context.Context, in *vo.DeptCreateReq) error {
	data := &v1.DeptCreateUpdateReq{
		DeptId:    0,
		ParentId:  in.ParentId,
		Ancestors: strings.TrimSpace(in.Ancestors),
		DeptName:  strings.TrimSpace(in.DeptName),
		Sort:      int64(in.Sort),
		Leader:    strings.TrimSpace(in.Leader),
		Phone:     strings.TrimSpace(in.Phone),
		Email:     strings.TrimSpace(in.Email),
		Status:    gconv.Int64(in.Status),
		CreateBy:  in.CreateBy,
	}
	_, err := r.data.sc.DeptCreate(ctx, data)
	return err
}

func (r *deptRepo) DeptUpdate(ctx context.Context, in *vo.DeptUpdateReq) error {
	data := &v1.DeptCreateUpdateReq{
		DeptId:    in.DeptId,
		ParentId:  in.ParentId,
		Ancestors: strings.TrimSpace(in.Ancestors),
		DeptName:  strings.TrimSpace(in.DeptName),
		Sort:      int64(in.Sort),
		Leader:    strings.TrimSpace(in.Leader),
		Phone:     strings.TrimSpace(in.Phone),
		Email:     strings.TrimSpace(in.Email),
		Status:    gconv.Int64(in.Status),
		UpdateBy:  in.UpdateBy,
	}
	_, err := r.data.sc.DeptUpdate(ctx, data)
	return err
}

func (r *deptRepo) DeptDelete(ctx context.Context, ids []int64) error {
	_, err := r.data.sc.DeptDelete(ctx, &v1.DeptDeleteReq{DeptIds: ids})
	return err
}
