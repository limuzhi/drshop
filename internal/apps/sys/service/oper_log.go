package service

import (
	"context"
	"drpshop/api/common"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/model"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/gogf/gf/util/gconv"
)

//操作日志列表
func (s *SysService) OperLogList(ctx context.Context, in *v1.OperLogListReq) (*v1.OperLogListRes, error) {
	relpy, total, err := s.operlogc.ListOperLog(ctx, in)
	if err != nil {
		return nil, errors.Unauthorized("OperLogError", "获取操作日志列表失败")
	}
	list := make([]*v1.OperLogInfoRes, len(relpy))
	for k, v := range relpy {
		info := s.operlogc.DtoOut(v)
		list[k] = info
	}
	return &v1.OperLogListRes{Total: total, List: list}, nil
}

//日志详细
func (s *SysService) OperLogInfo(ctx context.Context, in *v1.OperLogInfoReq) (*v1.OperLogInfoRes, error) {
	reply, err := s.operlogc.GetInfoById(ctx, in.OperId)
	if err != nil {
		return nil, errors.Unauthorized("OperLogError", "获取操作日志失败")
	}
	return s.operlogc.DtoOut(reply), nil
}

//操作日志删除
func (s *SysService) OperLogDelete(ctx context.Context, in *v1.OperLogDeleteReq) (*common.NullRes, error) {
	err := s.operlogc.BatchDeleteByIds(ctx, in.OperId)
	if err != nil {
		return nil, err
	}
	return &common.NullRes{}, nil
}

//操作日志清空
func (s *SysService) OperLogClear(ctx context.Context, _ *common.NullReq) (*common.NullRes, error) {
	err := s.operlogc.ClearOperLog(ctx)
	return &common.NullRes{}, err
}

//操作日志保存
func (s *SysService) OperLogSave(ctx context.Context, in *v1.OperLogSaveReq) (*common.NullRes, error) {
	if len(in.LogList) > 0 {
		paths := make([]string, 0)
		apiMap := make(map[int]methodPathApi)
		inserts := make([]*model.SysOperLog, len(in.LogList))
		for k, v := range in.LogList {
			inData := &model.SysOperLog{
				Title:         v.Title,
				BusinessType:  gconv.Int(v.BusinessType),
				Method:        v.Method,
				RequestMethod: v.RequestMethod,
				OperatorType:  gconv.Int(v.OperatorType),
				OperName:      v.OperName,
				OperUrl:       v.OperUrl,
				OperIp:        v.OperIp,
				OperLocation:  v.OperLocation,
				OperParam:     v.OperParam,
				JsonResult:    v.JsonResult,
				Status:        v.Status,
				ErrorMsg:      v.ErrorMsg,
				OperTime:      v.OperTime,
				TimeCost:      v.TimeCost,
				UserId:        v.UserId,
			}
			inserts[k] = inData
			info := methodPathApi{
				path:   v.Method,
				method: v.RequestMethod,
			}
			apiMap[k] = info
			paths = append(paths, v.Method)
		}
		apiList, _ := s.apic.GetListByPath(ctx, paths)
		if len(apiList) > 0 {
			for k, v := range apiMap {
				for _, val := range apiList {
					if val.Method == v.method && val.Path == v.path {
						inserts[k].Title = val.Title
						break
					}
				}
			}
		}
		if err := s.operlogc.BatchInsertOperLog(ctx, inserts); err != nil {
			return nil, err
		}
	}
	return &common.NullRes{}, nil
}

type methodPathApi struct {
	path   string
	method string
}
