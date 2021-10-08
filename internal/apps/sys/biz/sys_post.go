package biz

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/global"
	"drpshop/internal/apps/sys/model"
	"github.com/go-kratos/kratos/v2/log"
)

type SysPostRepo interface {
	//列表
	ListPost(ctx context.Context, in *v1.PostListReq) ([]*model.SysPost, int64, error)
	//创建
	CreatePost(ctx context.Context, in *model.SysPost) error
	//修改
	UpdatePost(ctx context.Context, id int64, in *model.SysPost) error
	//获取
	GetInfoById(ctx context.Context, id int64) (*model.SysPost, error)
	//删除
	BatchDeleteByIds(ctx context.Context, ids []int64) error
	//根据id获取列表
	ListPostByIds(ctx context.Context, ids []int64) ([]*model.SysPost, error)
	//获取所有列表选择
	AllListPost(ctx context.Context) ([]*model.SysPost, error)
}

type SysPostUsecase struct {
	repo SysPostRepo
	log  *log.Helper
}

func NewSysPostUsecase(repo SysPostRepo, logger log.Logger) *SysPostUsecase {
	return &SysPostUsecase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "sys/biz/sys_post")),
	}
}

//列表
func (uc *SysPostUsecase) ListPost(ctx context.Context, in *v1.PostListReq) ([]*model.SysPost, int64, error) {
	if in.PageInfo.PageNum <= 0 {
		in.PageInfo.PageNum = 1
	}
	if in.PageInfo.PageSize <= 0 {
		in.PageInfo.PageSize = 10
	}
	return uc.repo.ListPost(ctx, in)
}

//创建
func (uc *SysPostUsecase) CreatePost(ctx context.Context, in *model.SysPost) error {
	return uc.repo.CreatePost(ctx, in)
}

//修改
func (uc *SysPostUsecase) UpdatePost(ctx context.Context, id int64, in *model.SysPost) error {
	return uc.repo.UpdatePost(ctx, id, in)
}

//获取
func (uc *SysPostUsecase) GetInfoById(ctx context.Context, id int64) (*model.SysPost, error) {
	return uc.repo.GetInfoById(ctx, id)
}

//删除
func (uc *SysPostUsecase) BatchDeleteByIds(ctx context.Context, ids []int64) error {
	return uc.repo.BatchDeleteByIds(ctx, ids)
}

//根据id获取列表
func (uc *SysPostUsecase) ListPostByIds(ctx context.Context, ids []int64) ([]*model.SysPost, error) {
	return uc.repo.ListPostByIds(ctx, ids)
}

//获取所有列表选择
func (uc *SysPostUsecase) AllListPost(ctx context.Context) ([]*model.SysPost, error) {
	return uc.repo.AllListPost(ctx)
}

//===============
func (uc *SysPostUsecase) DtoOut(in *model.SysPost) *v1.PostDetailRes {
	out := &v1.PostDetailRes{
		PostId:    in.PostId,
		PostCode:  in.PostCode,
		PostName:  in.PostName,
		PostSort:  in.PostSort,
		Status:    in.Status,
		Remark:    in.Remark,
		CreateBy:  in.CreateBy,
		UpdateBy:  in.UpdateBy,
		CreatedAt: global.GetDateByUnix(in.CreatedAt),
		UpdatedAt: global.GetDateByUnix(in.UpdatedAt),
	}
	return out
}
