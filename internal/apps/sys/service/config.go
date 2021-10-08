package service

import (
	"context"
	"drpshop/api/common"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/model"
	"github.com/go-kratos/kratos/v2/errors"
	"strings"
)

func (s *SysService) ConfigInfoByKey(ctx context.Context, req *v1.ConfigInfoByKeyReq) (*v1.ConfigInfoByKeyRes, error) {
	configKey := strings.TrimSpace(req.ConfigKey)
	if configKey == "" {
		return nil, errors.Unauthorized("configKey", "configKey参数不能为空")
	}
	info, err := s.configc.GetInfoByConfigKey(ctx, req.ConfigKey)
	if err != nil {
		return nil, errors.Unauthorized("configKey", err.Error())
	}
	return &v1.ConfigInfoByKeyRes{
		ConfigId:    info.ConfigId,
		ConfigKey:   info.ConfigKey,
		ConfigValue: info.ConfigValue,
	}, nil
}

func (s *SysService) ConfigDetail(ctx context.Context, req *v1.ConfigDetailReq) (*v1.ConfigDetailRes, error) {
	reply, err := s.configc.GetInfoById(ctx, req.ConfigId)
	if err != nil {
		return nil, err
	}
	return s.configc.DtoOut(reply), nil
}
func (s *SysService) ConfigList(ctx context.Context, req *v1.ConfigListReq) (*v1.ConfigListRes, error) {
	reply, total, err := s.configc.ListConfig(ctx, req)
	if err != nil {
		return nil, err
	}
	out := make([]*v1.ConfigDetailRes, len(reply))
	for k, v := range reply {
		info := s.configc.DtoOut(v)
		out[k] = info
	}
	return &v1.ConfigListRes{List: out, Total: total}, err
}

func (s *SysService) ConfigCreate(ctx context.Context, req *v1.ConfigCreateUpdateReq) (*common.NullRes, error) {
	data := &model.SysConfig{
		ConfigName:  req.ConfigName,
		ConfigKey:   req.ConfigKey,
		ConfigValue: req.ConfigValue,
		ConfigType:  int(req.ConfigType),
		IsFrontend:  int(req.IsFrontend),
		Remark:      req.Remark,
	}
	data.CreateBy = req.CreateBy
	err := s.configc.CreateConfig(ctx, data)
	return &common.NullRes{}, err
}

func (s *SysService) ConfigUpdate(ctx context.Context, req *v1.ConfigCreateUpdateReq) (*common.NullRes, error) {
	data := &model.SysConfig{
		ConfigId:    req.ConfigId,
		ConfigName:  req.ConfigName,
		ConfigKey:   req.ConfigKey,
		ConfigValue: req.ConfigValue,
		ConfigType:  int(req.ConfigType),
		IsFrontend:  int(req.IsFrontend),
		Remark:      req.Remark,
	}
	data.UpdateBy = req.UpdateBy
	err := s.configc.UpdateConfig(ctx, req.ConfigId, data)
	return &common.NullRes{}, err
}
func (s *SysService) ConfigDelete(ctx context.Context, req *v1.ConfigDeleteReq) (*common.NullRes, error) {
	err := s.configc.BatchDeleteByIds(ctx, req.ConfigIds)
	return &common.NullRes{}, err
}
