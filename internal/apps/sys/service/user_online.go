package service

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/data/model"
	"github.com/go-kratos/kratos/v2/errors"
)

func (s *SysService) SaveUserOnline(ctx context.Context, req *v1.UserOnlineReq) (*v1.CommonRes, error) {
	err := s.userOnlinec.SaveUserOnline(ctx, &model.SysUserOnline{
		Username:   req.Username,
		UserId:     req.UserId,
		UUID:       req.Uuid,
		Token:      req.Token,
		Ip:         req.Ip,
		Explorer:   req.Explorer,
		Os:         req.Os,
		CreateTime: req.CreateTime,
	})
	if err != nil {
		return nil, errors.Unauthorized("userOnlineError", "用户在线状态保存")
	}
	return &v1.CommonRes{}, nil
}
