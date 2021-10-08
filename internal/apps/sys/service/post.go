package service

import (
	"context"
	"drpshop/api/common"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/model"
)

//岗位
func (s *SysService) PostDetail(ctx context.Context, in *v1.PostDetailReq) (*v1.PostDetailRes, error) {
	reply, err := s.postc.GetInfoById(ctx, in.PostId)
	if err != nil {
		return nil, err
	}
	return s.postc.DtoOut(reply), nil
}

func (s *SysService) PostList(ctx context.Context, in *v1.PostListReq) (*v1.PostListRes, error) {
	reply, total, err := s.postc.ListPost(ctx, in)
	if err != nil {
		return nil, err
	}
	out := make([]*v1.PostDetailRes, len(reply))
	for k, v := range reply {
		info := s.postc.DtoOut(v)
		out[k] = info
	}
	return &v1.PostListRes{List: out, Total: total}, err
}

func (s *SysService) PostCreate(ctx context.Context, in *v1.PostCreateUpdateReq) (*common.NullRes, error) {
	data := &model.SysPost{
		PostCode: in.PostCode,
		PostName: in.PostName,
		PostSort: in.PostSort,
		Status:   in.Status,
		Remark:   in.Remark,
	}
	data.CreateBy = in.CreateBy
	err := s.postc.CreatePost(ctx, data)
	return &common.NullRes{}, err
}

func (s *SysService) PostUpdate(ctx context.Context, in *v1.PostCreateUpdateReq) (*common.NullRes, error) {
	data := &model.SysPost{
		PostId:   in.PostId,
		PostCode: in.PostCode,
		PostName: in.PostName,
		PostSort: in.PostSort,
		Status:   in.Status,
		Remark:   in.Remark,
	}
	data.UpdateBy = in.UpdateBy
	err := s.postc.UpdatePost(ctx, in.PostId, data)
	return &common.NullRes{}, err
}
func (s *SysService) PostDelete(ctx context.Context, in *v1.PostDeleteReq) (*common.NullRes, error) {
	err := s.postc.BatchDeleteByIds(ctx, in.PostIds)
	return &common.NullRes{}, err
}
