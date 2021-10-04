package service

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/data/model"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/gogf/gf/util/gconv"
)

//获取字典分类 ---
//添加字典分类数据
func (s *SysService) DictTypeAdd(ctx context.Context, in *v1.DictTypeAddReq) (*v1.CommonRes, error) {
	data := &model.SysDictType{
		DictName: in.DictName,
		DictType: in.DictType,
		Remark:   in.Remark,
		Status:   gconv.Int(in.Status),
	}
	data.CreateBy = in.CreateBy
	err := s.dictTypec.CreateDictType(ctx, data)
	return &v1.CommonRes{}, err
}

//获取字典分类列表
func (s *SysService) DictTypeList(ctx context.Context, in *v1.DictTypeListReq) (*v1.DictTypeListRes, error) {
	if in.PageInfo.PageNum <= 0 {
		in.PageInfo.PageNum = 1
	}
	if in.PageInfo.PageSize <= 0 {
		in.PageInfo.PageSize= 10
	}
	reply, total, err := s.dictTypec.ListDictType(ctx, in)
	if err != nil {
		return nil, err
	}
	list := make([]*v1.DictTypeListData, len(reply))
	for k, v := range reply {
		info := s.dictTypec.DtoOut(v)
		list[k] = info
	}
	return &v1.DictTypeListRes{List: list, Total: total}, err
}

//获取字典分类info
func (s *SysService) DictTypeInfo(ctx context.Context, in *v1.DictTypeInfoReq) (*v1.DictTypeListData, error) {
	reply, err := s.dictTypec.GetInfoById(ctx, in.DictId)
	if err != nil {
		return nil, err
	}
	return s.dictTypec.DtoOut(reply), nil
}

//修改字典分类数据
func (s *SysService) DictTypeUpdate(ctx context.Context, in *v1.DictTypeUpdateReq) (*v1.CommonRes, error) {
	data := &model.SysDictType{
		DictId:   in.DictId,
		DictName: in.DictName,
		DictType: in.DictType,
		Remark:   in.Remark,
		Status:   gconv.Int(in.Status),
	}
	data.UpdateBy = in.UpdateBy
	err := s.dictTypec.UpdateDictType(ctx, in.DictId, data)
	return &v1.CommonRes{}, err
}

//删除字典分类数据
func (s *SysService) DictTypeDelete(ctx context.Context, in *v1.DictTypeDeleteReq) (*v1.CommonRes, error) {
	if len(in.DictIds) == 0 {
		return nil, errors.New(401, "参数不能为空", "参数不能为空")
	}
	err := s.dictTypec.BatchDeleteByIds(ctx, in.DictIds)
	return &v1.CommonRes{}, err
}

//获取字典选择框列表
func (s *SysService) DictTypeOptionSelect(ctx context.Context, _ *v1.CommonReq) (*v1.DictTypeListRes, error) {
	reply, err := s.dictTypec.GetAllDictType(ctx)
	if err != nil {
		return nil, err
	}
	list := make([]*v1.DictTypeListData, len(reply))
	for k, v := range reply {
		info := s.dictTypec.DtoOut(v)
		list[k] = info
	}
	return &v1.DictTypeListRes{List: list, Total: int64(len(list))}, err
}
