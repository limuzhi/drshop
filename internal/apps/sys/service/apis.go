package service

import (
	"context"
	"drpshop/api/common"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/model"
	"errors"
)

//接口
//接口详情
func (s *SysService) ApisInfo(ctx context.Context, in *v1.ApisInfoReq) (*v1.ApisInfoRes, error) {
	reply, err := s.apic.GetApisInfo(ctx, in.ApiId)
	if err != nil {
		return nil, err
	}
	return s.apic.DtoOut(reply), nil
}

//获取接口列表
func (s *SysService) ApisList(ctx context.Context, in *v1.ApisListReq) (*v1.ApisListRes, error) {
	reply, total, err := s.apic.ListApis(ctx, in)
	if err != nil {
		return nil, err
	}
	list := make([]*v1.ApisInfoRes, len(reply))
	for k, v := range reply {
		info := s.apic.DtoOut(v)
		list[k] = info
	}
	return &v1.ApisListRes{List: list, Total: total}, nil
}

//获取接口树(按接口Category字段分类)
func (s *SysService) ApisTreeList(ctx context.Context, _ *common.NullReq) (*v1.ApisTreeListRes, error) {
	list, err := s.apic.GetApisListTree(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.ApisTreeListRes{List: list}, nil
}

//创建接口
func (s *SysService) ApisAdd(ctx context.Context, in *v1.ApisAddReq) (*common.NullRes, error) {
	insert := &model.SysApis{
		Handle:     in.Handle,
		Title:      in.Title,
		Path:       in.Path,
		Method:     in.Method,
		Category:   in.Category,
		Permission: in.Permission,
	}
	insert.CreateBy = in.CreateBy
	err := s.apic.CreateApis(ctx, insert)
	return &common.NullRes{}, err
}

//更新接口
func (s *SysService) ApisUpdate(ctx context.Context, in *v1.ApisUpdateReq) (*common.NullRes, error) {
	if in.ApiId < 1 {
		return nil, errors.New("apiId必须大于0")
	}
	data := &model.SysApis{
		ApiId:      in.ApiId,
		Handle:     in.Handle,
		Title:      in.Title,
		Path:       in.Path,
		Method:     in.Method,
		Category:   in.Category,
		Permission: in.Permission,
	}
	data.UpdateBy = in.UpdateBy
	err := s.apic.UpdateApis(ctx, in.ApiId, data)
	return &common.NullRes{}, err
}

//批量删除接口
func (s *SysService) ApisDelete(ctx context.Context, in *v1.ApisDeleteReq) (*common.NullRes, error) {
	if len(in.ApiIds) == 0 {
		return nil, errors.New("不能为空")
	}
	err := s.apic.BatchDeleteByIds(ctx, in.ApiIds)
	return &common.NullRes{}, err
}
