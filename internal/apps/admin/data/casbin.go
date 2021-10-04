package data

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/admin/biz"
	"github.com/go-kratos/kratos/v2/log"
	"sync"
)

type casbinRepo struct {
	data *Data
	log  *log.Helper
}

func NewCasbinRepo(data *Data, logger log.Logger) biz.CasbinRepo {
	return &casbinRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "admin/data/casbin")),
	}
}

var checkLock sync.Mutex

func (r *casbinRepo) CheckCasbin(ctx context.Context, userId int64, obj, act string) bool {
	checkLock.Lock()
	defer checkLock.Unlock()
	reply, err := r.data.sc.CheckCasbin(ctx, &v1.CheckCasbinReq{
		UserId: userId,
		Obj:    obj,
		Act:    act,
	})
	if err != nil {
		r.log.Info("CheckCasbin---err:",err.Error())
		return false
	}
	if reply.Pong != "ok" {
		return false
	}
	return true
}
