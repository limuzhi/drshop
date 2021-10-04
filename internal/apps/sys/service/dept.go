package service

import (
	"context"
	"drpshop/api/common"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/data/model"
	"errors"
	"fmt"
)

//部门
func (s *SysService) DeptDetail(ctx context.Context, req *v1.DeptDetailReq) (*v1.DeptDetailRes, error) {
	reply, err := s.deptc.GetInfoById(ctx, req.DeptId)
	if err != nil {
		return nil, err
	}
	return s.deptc.DtoOut(reply), nil
}

func (s *SysService) DeptList(ctx context.Context, req *v1.DeptListReq) (*v1.DeptListRes, error) {
	reply, total, err := s.deptc.ListDept(ctx, req)
	if err != nil {
		return nil, err
	}
	out := make([]*v1.DeptDetailRes, len(reply))
	for k, v := range reply {
		info := s.deptc.DtoOut(v)
		out[k] = info
	}
	return &v1.DeptListRes{List: out, Total: total}, err
}

func (s *SysService) DeptTree(ctx context.Context, req *v1.DeptListReq) (*v1.DeptListRes, error) {
	reply, total, err := s.deptc.ListDept(ctx, req)
	if err != nil {
		return nil, err
	}
	out := s.deptc.GetDeptListTree(0, reply)
	return &v1.DeptListRes{List: out, Total: total}, err
}

func (s *SysService) DeptCreate(ctx context.Context, req *v1.DeptCreateUpdateReq) (*common.NullRes, error) {
	data := &model.SysDept{
		ParentId:  req.ParentId,
		DeptName:  req.DeptName,
		Sort:      int(req.Sort),
		Leader:    req.Leader,
		Phone:     req.Phone,
		Email:     req.Email,
		Status:    int(req.Status),
	}
	data.CreateBy = req.CreateBy
	data.Ancestors = "0"
	if req.ParentId > 0 {
		info,err := s.deptc.GetInfoById(ctx,req.ParentId)
		if err != nil {
			return nil, errors.New("获取上级部门失败"+err.Error())
		}
		data.Ancestors = fmt.Sprintf("%s,%d",info.Ancestors,info.DeptId)
	}
	err := s.deptc.CreateDept(ctx, data)
	return &common.NullRes{}, err
}

func (s *SysService) DeptUpdate(ctx context.Context, req *v1.DeptCreateUpdateReq) (*common.NullRes, error) {
	data := &model.SysDept{
		DeptId:    req.DeptId,
		ParentId:  req.ParentId,
		Ancestors: req.Ancestors,
		DeptName:  req.DeptName,
		Sort:      int(req.Sort),
		Leader:    req.Leader,
		Phone:     req.Phone,
		Email:     req.Email,
		Status:    int(req.Status),
	}
	data.UpdateBy = req.UpdateBy
	data.Ancestors = "0"
	if req.ParentId > 0 {
		info,err := s.deptc.GetInfoById(ctx,req.ParentId)
		if err != nil {
			return nil, errors.New("获取上级部门失败"+err.Error())
		}
		data.Ancestors = fmt.Sprintf("%s,%d",info.Ancestors,info.DeptId)
	}
	err := s.deptc.UpdateDept(ctx, req.DeptId, data)
	return &common.NullRes{}, err
}

func (s *SysService) DeptDelete(ctx context.Context, req *v1.DeptDeleteReq) (*common.NullRes, error) {
	err := s.deptc.BatchDeleteByIds(ctx, req.DeptIds)
	return &common.NullRes{}, err
}
