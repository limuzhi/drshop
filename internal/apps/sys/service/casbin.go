package service

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/data/model"
	"drpshop/pkg/token"
	"github.com/go-kratos/kratos/v2/errors"
	"strings"
)

//检测rbac
func (s *SysService) CheckCasbin(ctx context.Context, req *v1.CheckCasbinReq) (*v1.CommonRes, error) {
	userId := token.FormGlobalUidContext(ctx)
	if userId <= 0 || req.UserId != userId {
		return nil, errors.New(401, "CheckCasbin", "用户未登录")
	}
	reply, err := s.userc.GetInfoById(ctx, userId)
	if err != nil {
		return nil, err
	}
	if reply.Status != int64(model.Enabled) {
		return nil, errors.New(401, "CheckCasbin", "当前用户已被禁用")
	}
	var subs []string
	for _, v := range reply.Roles {
		if v.Status == int64(model.Enabled) {
			subs = append(subs, v.RoleKey)
		}
	}
	if len(subs) == 0 {
		return nil, errors.New(401, "CheckCasbin", "当前用户无角色")
	}
	obj := strings.TrimSpace(req.Obj)
	act := strings.TrimSpace(req.Act)
	isPass := s.rolec.CheckCasbinRole(subs, obj, act)
	if !isPass {
		return nil, errors.New(401, "CheckCasbin", "没有权限")
	}
	return &v1.CommonRes{Pong: "ok"}, nil
}
