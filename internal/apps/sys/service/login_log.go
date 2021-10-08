package service

import (
	"context"
	"drpshop/api/common"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/model"
	"drpshop/pkg/util"
	"github.com/go-kratos/kratos/v2/errors"
)

func (s *SysService) SaveLoginlog(ctx context.Context, req *v1.LoginLogReq) (*v1.CommonRes, error) {

	err := s.loginLogc.SaveLoginlog(ctx, &model.SysLoginLog{
		LoginName:     req.LoginName,
		LoginUid:      req.LoginUid,
		Ipaddr:        req.Ipaddr,
		LoginLocation: util.GetCityByIp(req.Ipaddr),
		Browser:       req.Browser,
		Os:            req.Os,
		Msg:           req.Msg,
		Status:        req.Status,
		LoginTime:     req.LoginTime,
	})
	if err != nil {
		return nil, errors.Unauthorized("LoginLogError", "插入登录日志失败")
	}
	return &v1.CommonRes{}, nil
}

//登录日志列表
func (s *SysService) LoginlogList(ctx context.Context, in *v1.LoginlogListReq) (*v1.LoginlogListRes, error) {
	if in.PageInfo.PageNum <= 0 {
		in.PageInfo.PageNum = 1
	}
	if in.PageInfo.PageSize <= 0 {
		in.PageInfo.PageSize = 10
	}
	reply, total, err := s.loginLogc.ListLoginLog(ctx, in)
	if err != nil {
		return nil, err
	}
	list := s.loginLogc.DtoOut(reply)
	return &v1.LoginlogListRes{List: list, Total: total}, nil
}

//删除
func (s *SysService) LoginlogDelete(ctx context.Context, in *v1.LoginlogDeleteReq) (*v1.CommonRes, error) {
	err := s.loginLogc.BatchDeleteByIds(ctx, in.LoginIds)
	return &v1.CommonRes{}, err
}

//清空
func (s *SysService) LoginlogClear(ctx context.Context, _ *common.NullReq) (*v1.CommonRes, error) {
	err := s.loginLogc.ClearLoginLog(ctx)
	return &v1.CommonRes{}, err
}
