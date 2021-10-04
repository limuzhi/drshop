package data

import (
	"context"
	"drpshop/api/common"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/admin/biz"
	"drpshop/internal/apps/admin/vo"
	"github.com/go-kratos/kratos/v2/log"
	"strings"
)

type loginLogRepo struct {
	data *Data
	log  *log.Helper
}

func NewLoginLogRepo(data *Data, logger log.Logger) biz.LoginLogRepo {
	return &loginLogRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "admin/data/login_log")),
	}
}

func (r *loginLogRepo) List(ctx context.Context, in *vo.LoginLogSearchReq) (*v1.LoginlogListRes, error) {
	params := &v1.LoginlogListReq{
		LoginName:     strings.TrimSpace(in.LoginName),
		Ipaddr:        strings.TrimSpace(in.Ipaddr),
		LoginLocation: strings.TrimSpace(in.LoginLocation),
		Status:        in.Status,
		PageInfo: &common.PageReq{
			PageNum:   int64(in.GetPageIndex()),
			PageSize:  int64(in.GetPageSize()),
			BeginTime: in.GetBeginTime(),
			EndTime:   in.GetEndTime(),
		},
	}
	return r.data.sc.LoginlogList(ctx, params)
}

func (r *loginLogRepo) BatchDelete(ctx context.Context, ids []int64) error {
	_, err := r.data.sc.LoginlogDelete(ctx, &v1.LoginlogDeleteReq{LoginIds: ids})
	return err
}

func (r *loginLogRepo) Clear(ctx context.Context) error {
	_, err := r.data.sc.LoginlogClear(ctx, &common.NullReq{})
	return err
}
