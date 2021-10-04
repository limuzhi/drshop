package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/util/gconv"
	"gorm.io/gorm"
	"strings"
	"time"

	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/biz"
	"drpshop/internal/apps/sys/data/model"
	"github.com/go-kratos/kratos/v2/log"
)

type sysUserRepo struct {
	data *Data
	log  *log.Helper
}

func NewSysUserRepo(data *Data, logger log.Logger) biz.SysUserRepo {
	return &sysUserRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "sys/data/sys_user")),
	}
}

//用户登录
func (r *sysUserRepo) GetUsersForAccounts(ctx context.Context, account string) (*model.SysUser, error) {
	var info *model.SysUser
	err := r.data.db.WithContext(ctx).Model(&model.SysUser{}).
		Where("username = ? OR mobile = ? ", account, account).
		Preload("Roles").First(&info).Error
	return info, err
}

//获取用户
func (r *sysUserRepo) GetInfoById(ctx context.Context, uid int64) (*model.SysUser, error) {
	var firstUser model.SysUser
	err := r.data.db.WithContext(ctx).Where("user_id = ?", uid).
		Preload("Roles").First(&firstUser).Error
	return &firstUser, err
}

//创建用户 gorm 文档地址 https://gorm.io/zh_CN/docs/associations.html#delete_with_select
func (r *sysUserRepo) CreateUser(ctx context.Context, in *model.SysUser) error {
	return r.data.db.WithContext(ctx).Omit("Roles.*").Omit("Posts.*").Create(in).Error
}

//更新用户
func (r *sysUserRepo) UpdateUser(ctx context.Context, in *model.SysUser) error {
	var count int64
	r.data.db.WithContext(ctx).Model(&model.SysUser{}).
		Where("user_id != ? AND mobile = ? ", in.UserId, in.Mobile).Count(&count)
	if count > 0 {
		return errors.New("手机号已存在")
	}
	// 事务处理， tx 处理数据
	tx := r.data.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}

	err := tx.Model(in).Where("user_id = ?", in.UserId).
		Omit("Roles.*").Omit("Posts.*").Updates(in).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Model(in).Omit("Roles.*").Association("Roles").Replace(in.Roles)
	// 如果更新成功就更新用户信息缓存
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Model(in).Omit("Posts.*").Association("Posts").Replace(in.Posts)
	// 如果更新成功就更新用户信息缓存
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

//列表用户信息
func (r *sysUserRepo) ListUser(ctx context.Context, in *v1.UserListReq) ([]*model.SysUser, int64, error) {
	var list []*model.SysUser
	table := r.data.db.WithContext(ctx)
	table = table.Model(&model.SysUser{}).Order("created_at DESC")
	username := strings.TrimSpace(in.Username)
	if username != "" {
		keywords := fmt.Sprintf("%%%s%%", username)
		table = table.Where("username LIKE ? OR nickname LIKE ?", keywords, keywords)
	}
	mobile := strings.TrimSpace(in.Mobile)
	if mobile != "" {
		table = table.Where("mobile LIKE ?", fmt.Sprintf("%%%s%%", mobile))
	}
	email := strings.TrimSpace(in.Email)
	if email != "" {
		table = table.Where("email LIKE ?", fmt.Sprintf("%%%s%%", email))
	}
	if in.Status != "" {
		table = table.Where("status = ?", gconv.Int(in.Status))
	}
	if in.StartTime > 0 {
		table = table.Where("created_at >= ?", in.StartTime)
	}
	if in.EndTime > 0 {
		table = table.Where("created_at <= ?", in.EndTime)
	}
	// 当pageNum > 0 且 pageSize > 0 才分页
	var total int64
	err := table.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	if in.PageNum <= 0 {
		in.PageNum = 1
	}
	pageNum := int(in.PageNum)
	pageSize := int(in.PageSize)
	err = table.Offset((pageNum - 1) * pageSize).Limit(pageSize).
		Preload("Roles").
		Preload("Posts").Find(&list).Error
	return list, total, err
}

//批量删除用户数据
func (r *sysUserRepo) BatchDeleteUserByIds(ctx context.Context, ids []int64) error {
	var list []*model.SysUser
	tx := r.data.db.WithContext(ctx)
	err := tx.Model(&model.SysUser{}).Where("user_id IN (?)", ids).Find(&list).Error
	if err != nil {
		return err
	}
	err = tx.Model(&model.SysUser{}).Select("Roles", "Posts").Unscoped().Delete(&list).Error

	return err
}

//更新密码
func (r *sysUserRepo) ChangePwd(ctx context.Context, uid int64, salt, hashNewPasswd string) error {
	err := r.data.db.WithContext(ctx).Model(&model.SysUser{}).
		Where("user_id = ?", uid).Updates(map[string]interface{}{
		"salt":     salt,
		"password": hashNewPasswd,
	}).Error
	// 如果更新密码成功，则更新当前用户信息缓存  TODO
	return err
}

//修改用户状态
func (r *sysUserRepo) UpdateUserStatus(ctx context.Context, uid, status int64) error {
	return r.data.db.WithContext(ctx).Model(&model.SysUser{}).Where("user_id = ?", uid).
		Updates(map[string]interface{}{"status": status}).Error
}

func (r *sysUserRepo) GetPermissions(roleIds []int64) ([]string, error) {
	err := r.data.casbinEnforcer.LoadPolicy()
	if err != nil {
		return nil, err
	}
	userButtons := make([]string, 0)
	for _, roleId := range roleIds {
		//查询当前权限
		gp := r.data.casbinEnforcer.GetFilteredPolicy(0, gconv.String(roleId))
		for _, p := range gp {
			userButtons = append(userButtons, p[1])
		}
	}
	return userButtons, nil
}

//废弃掉了
func (r *sysUserRepo) CreateUserPostRole(ctx context.Context, in *model.SysUser,
	postIds []int64, roleIds []int64) error {
	// 事务处理， tx 处理数据
	tx := r.data.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	var count int64
	tx.Model(&model.SysUser{}).Where("username = ?", in.Username).Count(&count)
	if count > 0 {
		tx.Rollback()
		return errors.New("用户名已经存在")
	}
	tx.Model(&model.SysUser{}).Where("mobile = ?", in.Mobile).Count(&count)
	if count > 0 {
		tx.Rollback()
		return errors.New("手机号已经存在")
	}
	if err := tx.Model(&model.SysUser{}).Create(in).Error; err != nil {
		tx.Rollback()
		return err
	}
	_ = tx.Model(&model.SysUserRoles{}).Where("user_id = ?", in.UserId).
		Unscoped().Delete(&model.SysUserRoles{}).Error
	if len(roleIds) > 0 {
		roleList := make([]*model.SysUserRoles, 0)
		for _, v := range roleIds {
			roleList = append(roleList, &model.SysUserRoles{
				UserId: in.UserId,
				RoleId: v,
			})
		}
		if err := tx.Model(&model.SysUserRoles{}).Create(roleList).Error; err != nil {
			tx.Rollback()
			return err
		}

	}
	////删除旧岗位信息
	_ = tx.Model(&model.SysUserPost{}).Where("user_id = ?", in.UserId).
		Unscoped().Delete(&model.SysUserPost{}).Error
	if len(postIds) > 0 {
		postList := make([]*model.SysUserPost, 0)
		for _, v := range postIds {
			postList = append(postList, &model.SysUserPost{
				UserId: in.UserId,
				PostId: v,
			})
		}
		if err := tx.Model(&model.SysUserPost{}).Create(postList).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (r *sysUserRepo) LoginUpdateUser(ctx context.Context, userId, loginTime int64, ip string) error {
	return r.data.db.WithContext(ctx).Model(&model.SysUser{}).Where("user_id = ?", userId).
		Updates(map[string]interface{}{
			"last_login_ip":   ip,
			"last_login_time": loginTime,
			"login_count":     gorm.Expr("login_count + ?", 1),
		}).Error
}

//设置用户redis的TOKEN缓存--暂不用
func (r *sysUserRepo) SetTokenCache(ctx context.Context, cacheKey string,
	data map[string]interface{}, expira int64) error {
	return r.data.rdb.SetEX(ctx, cacheKey, data, time.Second*time.Duration(expira)).Err()
}

func (r *sysUserRepo) GetUserProfile(ctx context.Context, uid int64) (*model.SysUser, error) {
	var info *model.SysUser
	err := r.data.db.WithContext(ctx).Model(&model.SysUser{}).
		Where("user_id = ?", uid).Preload("Roles", "status = ?", model.Enabled).
		Preload("Posts", "status = ?", model.Enabled).First(&info).Error
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (r *sysUserRepo) UpdateUserProfile(ctx context.Context, uid int64,data map[string]interface{}) error {
	tx := r.data.db.WithContext(ctx).Model(&model.SysUser{}).Where("user_id = ?",uid)
	return  tx.Updates(data).Error
}
