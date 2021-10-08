package service

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/model"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/gogf/gf/util/gconv"
)

//字典数据---
//添加字典数据
func (s *SysService) DictDataAdd(ctx context.Context, in *v1.DictDataAddReq) (*v1.CommonRes, error) {
	data := &model.SysDictData{
		DictSort:  gconv.Int(in.DictSort),
		DictLabel: in.DictLabel,
		DictValue: in.DictValue,
		DictType:  in.DictType,
		CssClass:  in.CssClass,
		ListClass: in.ListClass,
		Remark:    in.Remark,
		IsDefault: gconv.Int(in.IsDefault),
		Status:    gconv.Int(in.Status),
	}
	data.CreateBy = in.CreateBy
	err := s.dictDatac.CreateDictData(ctx, data)
	return &v1.CommonRes{}, err
}

//字典数据列表
func (s *SysService) DictDataList(ctx context.Context, in *v1.DictDataListReq) (*v1.DictDataListRes, error) {
	if in.PageInfo.PageNum <= 0 {
		in.PageInfo.PageNum = 1
	}
	if in.PageInfo.PageSize <= 0 {
		in.PageInfo.PageSize = 10
	}
	reply, total, err := s.dictDatac.ListDictData(ctx, in)
	if err != nil {
		return nil, err
	}
	list := make([]*v1.DictDataListData, len(reply))
	for k, v := range reply {
		info := s.dictDatac.DtoOut(v)
		list[k] = info
	}
	return &v1.DictDataListRes{List: list, Total: total}, nil
}

//修改字典数据
func (s *SysService) DictDataUpdate(ctx context.Context, in *v1.DictDataUpdateReq) (*v1.CommonRes, error) {
	data := &model.SysDictData{
		DictCode:  in.DictCode,
		DictSort:  gconv.Int(in.DictSort),
		DictLabel: in.DictLabel,
		DictValue: in.DictValue,
		DictType:  in.DictType,
		CssClass:  in.CssClass,
		ListClass: in.ListClass,
		Remark:    in.Remark,
		IsDefault: gconv.Int(in.IsDefault),
		Status:    gconv.Int(in.Status),
	}
	data.UpdateBy = in.UpdateBy
	err := s.dictDatac.UpdateDictData(ctx, in.DictCode, data)
	return &v1.CommonRes{}, err
}

//删除字典数据
func (s *SysService) DictDataDelete(ctx context.Context, in *v1.DictDataDeleteReq) (*v1.CommonRes, error) {
	if len(in.DictCodes) == 0 {
		return nil, errors.Unauthorized("DictDeleteError", "参数不能为空")
	}
	err := s.dictDatac.BatchDeleteByIds(ctx, in.DictCodes)
	return &v1.CommonRes{}, err
}

//通过编码获取字典数据
func (s *SysService) DictDataInfoByDictCode(ctx context.Context, in *v1.DictDataInfoByDictCodeReq) (*v1.DictDataListData, error) {
	reply, err := s.dictDatac.GetInfoById(ctx, in.DictCode)
	if err != nil {
		return nil, err
	}
	return s.dictDatac.DtoOut(reply), err
}

//数据字典根据key获取
func (s *SysService) DictDataListByDictType(ctx context.Context, req *v1.DictDataListByDictTypeReq) (*v1.DictDataListRes, error) {
	res, err := s.dictDatac.ListByDictType(ctx, req.DictType)
	if err != nil {
		return nil, err
	}
	list := make([]*v1.DictDataListData, 0)
	for _, v := range res {
		if v.Status != int(model.Enabled) {
			continue
		}
		info := s.dictDatac.DtoOut(v)
		list = append(list, info)
	}
	return &v1.DictDataListRes{List: list}, nil
}
