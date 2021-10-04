package biz

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/data/model"
	"drpshop/internal/apps/sys/global"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/util/gconv"
)

type SysDeptRepo interface {
	//列表
	ListDept(ctx context.Context, in *v1.DeptListReq) ([]*model.SysDept, int64, error)
	//删除
	BatchDeleteByIds(ctx context.Context, ids []int64) error
	//创建
	CreateDept(ctx context.Context, in *model.SysDept) error
	//修改
	UpdateDept(ctx context.Context, id int64, in *model.SysDept) error
	//获取
	GetInfoById(ctx context.Context, id int64) (*model.SysDept, error)
	//根据角色ID获取部门
	GetRoleDepts(ctx context.Context, roleId int64) ([]int64, error)
	//获取部门选择
	AllListDept(ctx context.Context) ([]*model.SysDept, error)
}

type SysDeptUsecase struct {
	repo SysDeptRepo
	log  *log.Helper
}

func NewSysDeptUsecase(repo SysDeptRepo, logger log.Logger) *SysDeptUsecase {
	return &SysDeptUsecase{repo: repo, log: log.NewHelper(log.With(logger, "module", "sys/biz/sys_dept"))}
}

//列表
func (uc *SysDeptUsecase) ListDept(ctx context.Context, in *v1.DeptListReq) ([]*model.SysDept, int64, error) {
	return uc.repo.ListDept(ctx, in)
}

//删除
func (uc *SysDeptUsecase) BatchDeleteByIds(ctx context.Context, ids []int64) error {
	return uc.repo.BatchDeleteByIds(ctx, ids)
}

//创建
func (uc *SysDeptUsecase) CreateDept(ctx context.Context, in *model.SysDept) error {
	return uc.repo.CreateDept(ctx, in)
}

//修改
func (uc *SysDeptUsecase) UpdateDept(ctx context.Context, id int64, in *model.SysDept) error {
	return uc.repo.UpdateDept(ctx,id, in)
}

//获取
func (uc *SysDeptUsecase) GetInfoById(ctx context.Context, id int64) (*model.SysDept, error) {
	return uc.repo.GetInfoById(ctx, id)
}

//根据角色ID获取部门
func (uc *SysDeptUsecase) GetRoleDepts(ctx context.Context, roleId int64) ([]int64, error) {
	return uc.repo.GetRoleDepts(ctx, roleId)
}

//获取部门选择
func (uc *SysDeptUsecase) AllListDept(ctx context.Context) ([]*model.SysDept, error) {
	list, _, err := uc.repo.ListDept(ctx, &v1.DeptListReq{Status: gconv.String(model.Enabled)})
	return list, err
}

//TODO==============

func (uc *SysDeptUsecase) GetDeptListTree(pid int64, list []*model.SysDept) []*v1.DeptDetailRes {
	tree := make([]*v1.DeptDetailRes, 0, len(list))
	for _, v := range list {
		info := &v1.DeptDetailRes{
			DeptId:    v.DeptId,
			ParentId:  v.ParentId,
			Ancestors: v.Ancestors,
			DeptName:  v.DeptName,
			Sort:      int64(v.Sort),
			Leader:    v.Leader,
			Phone:     v.Phone,
			Email:     v.Email,
			Status:    int64(v.Status),
			CreateBy:  v.CreateBy,
			UpdateBy:  v.UpdateBy,
			CreatedAt: global.GetDateByUnix(v.CreatedAt),
			UpdatedAt: global.GetDateByUnix(v.UpdatedAt),
		}
		if v.ParentId == pid {
			child := uc.GetDeptListTree(v.DeptId, list)
			if child != nil && len(child) > 0 {
				info.Children = child
			}
			tree = append(tree, info)
		}
	}
	return tree
}

func (uc *SysDeptUsecase) AllDeptList(ctx context.Context) ([]*model.SysDept, int64, error) {
	return uc.repo.ListDept(ctx, &v1.DeptListReq{})
}

func (uc *SysDeptUsecase) DtoOut(data *model.SysDept) *v1.DeptDetailRes {
	info := &v1.DeptDetailRes{
		DeptId:    data.DeptId,
		ParentId:  data.ParentId,
		Ancestors: data.Ancestors,
		DeptName:  data.DeptName,
		Sort:      int64(data.Sort),
		Leader:    data.Leader,
		Phone:     data.Phone,
		Email:     data.Email,
		Status:    int64(data.Status),
		CreateBy:  data.CreateBy,
		UpdateBy:  data.UpdateBy,
		CreatedAt: global.GetDateByUnix(data.CreatedAt),
		UpdatedAt: global.GetDateByUnix(data.UpdatedAt),
	}
	return info
}
