package data

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/admin/biz"
	"drpshop/internal/apps/admin/vo"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/util/gconv"
	"strings"
)

type roleRepo struct {
	data *Data
	log  *log.Helper
}

func NewRoleRepo(data *Data, logger log.Logger) biz.RoleRepo {
	return &roleRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "admin/data/role")),
	}
}

func (r *roleRepo) SearchRoleList(ctx context.Context, in *v1.RoleListReq) (*v1.RoleListRes, error) {
	return r.data.sc.RoleList(ctx, in)
}

func (r *roleRepo) InsertRole(ctx context.Context, in *vo.RoleAddReq) error {
	_, err := r.data.sc.RoleAdd(ctx, &v1.RoleAddReq{
		Sort:      in.Sort,
		Pid:       0,
		Name:      strings.TrimSpace(in.Name),
		RoleKey:   strings.TrimSpace(in.RoleKey),
		Remark:    in.Remark,
		Status:    gconv.Int64(in.Status),
		DataScope: 0, //数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限）
	})
	return err
}

func (r *roleRepo) EditRole(ctx context.Context, in *vo.RoleUpdateReq) error {
	_, err := r.data.sc.RoleUpdate(ctx, &v1.RoleUpdateReq{
		RoleId:  in.RoleId,
		Sort:    in.Sort,
		Pid:     0,
		Name:    strings.TrimSpace(in.Name),
		RoleKey: strings.TrimSpace(in.RoleKey),
		Remark:  in.Remark,
		Status:  gconv.Int64(in.Status),
	})
	return err
}

func (r *roleRepo) BatchDeleteRole(ctx context.Context, ids []int64) error {
	_, err := r.data.sc.RoleDelete(ctx, &v1.RoleDeleteReq{RoleIds: ids})
	return err
}

func (r *roleRepo) ChangeRoleStatus(ctx context.Context, id, status int64) error {
	_, err := r.data.sc.UpdateRoleStatus(ctx, &v1.UpdateRoleStatusReq{
		RoleId: id,
		Status: status,
	})
	return err
}

func (r *roleRepo) RoleMeunsByRoleIds(ctx context.Context, roleId int64) (*v1.QueryMenuByRoleIdRes, error) {
	return r.data.sc.GetMenusByRoleId(ctx, &v1.QueryMenuByRoleIdReq{
		RoleId: roleId,
	})
}

func (r *roleRepo) RoleApisByRoleIds(ctx context.Context, roleId int64) (*v1.QueryApisByRoleIdRes, error) {
	return r.data.sc.GetApisByRoleId(ctx, &v1.QueryApisByRoleIdReq{RoleId: roleId})
}

func (r *roleRepo) RoleMenusUpdate(ctx context.Context, userId, roleId int64, menuIds []int64) error {
	_, err := r.data.sc.UpdateMenuRole(ctx, &v1.UpdateMenuRoleReq{
		UserId:  userId,
		RoleId:  roleId,
		MenuIds: menuIds,
	})
	return err
}
func (r *roleRepo) RoleApisUpdate(ctx context.Context, userId, roleId int64, apiIds []int64) error {
	_, err := r.data.sc.UpdateRoleApisById(ctx, &v1.UpdateApisRoleReq{
		UserId: userId,
		RoleId: roleId,
		ApiIds: apiIds,
	})
	return err
}

