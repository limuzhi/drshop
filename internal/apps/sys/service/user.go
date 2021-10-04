package service

import (
	"context"
	"drpshop/api/common"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/data/model"
	"drpshop/internal/apps/sys/global"
	"drpshop/pkg/captcha"
	"drpshop/pkg/token"
	"drpshop/pkg/util"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"log"
	"math/rand"
	"strings"
	"time"
)

func (s *SysService) Login(ctx context.Context, req *v1.LoginReq) (*v1.LoginRes, error) {
	if req.Password == "" || req.Account == "" {
		return nil, errors.Unauthorized("LoginError", "账号密码不能为空")
	}
	if req.Captcha == "" {
		return nil, errors.Unauthorized("LoginError", "验证码不能为空")
	}
	if !s.store.Verify(req.CaptchaId, req.Captcha, true) {
		return nil, errors.Unauthorized("LoginError", "验证码错误")
	}
	tk, etime, err := s.userc.Login(ctx, req.Account, req.Password)
	if err != nil {
		return nil, errors.Unauthorized("LoginError", err.Error())
	}
	return &v1.LoginRes{
		Token:  tk,
		Expire: etime,
	}, nil
}

func (s *SysService) CaptchaImg(_ context.Context, req *v1.CaptchaImgReq) (*v1.CaptchaImgRes, error) {
	if id, b64s, err := captcha.DriverStringFunc(s.store, int(req.Height),
		int(req.Width), int(req.Length)); err != nil {
		return nil, errors.Unauthorized("验证码获取失败", err.Error())
	} else {
		return &v1.CaptchaImgRes{
			CaptchaId: id,
			B64S:      b64s,
		}, nil
	}
}

func (s *SysService) Captcha(_ context.Context, req *v1.CaptchaReq) (*common.NullRes, error) {
	if req.Captcha == "" {
		return nil, errors.Unauthorized("LoginError", "验证码不能为空")
	}
	if req.To == "" {
		return nil, errors.Unauthorized("LoginError", "接收者不能为空")
	}

	if !s.store.Verify(req.CaptchaId, req.Captcha, true) {
		return nil, errors.Unauthorized("LoginError", "验证码错误")
	}
	code := fmt.Sprint(rand.Intn(9999) + 1000)
	s.store.Set("captcha:"+req.To, code)
	if util.CheckEmail(req.To) {
		// todo 发送邮箱
	} else if util.CheckMobile(req.To) {
		// todo 发送短信
	}
	return &common.NullRes{}, nil
}

//获取用户
func (s *SysService) UserInfo(ctx context.Context, req *v1.UserInfoReq) (*v1.UserInfoRes, error) {
	uid := token.FormGlobalUidContext(ctx)
	if uid != req.UserId {
		return nil, errors.Unauthorized("UserInfoError", "用户信息失败")
	}
	info, err := s.userc.UserInfo(ctx, req.UserId)
	if err != nil {
		return nil, errors.Unauthorized("UserInfoError", "获取用户信息失败")
	}
	if info.Status != 2 {
		return nil, errors.Unauthorized("UserInfoError", "用户已被禁用")
	}

	roleKeys := make([]string, 0)
	roleIds := make([]int64, 0)
	isSuperAdmin := false
	for _, v := range info.Roles {
		if v.Status != 2 {
			continue
		}
		roleKeys = append(roleKeys, v.RoleKey)
		roleIds = append(roleIds, v.RoleId)
		if !isSuperAdmin && v.RoleKey == "admin" {
			isSuperAdmin = true
		}
	}
	var permissions []string
	if isSuperAdmin {
		permissions = []string{"*:*:*"}
	} else {
		for _, v := range roleKeys {
			permissionsList, apiErr := s.rolec.GetPermissionsByRoleKey(v)
			if apiErr != nil {
				continue
			}
			if len(permissionsList) > 0 {
				permissions = append(permissions, permissionsList...)
			}
		}
	}
	return &v1.UserInfoRes{
		Avatar:       info.Avatar,
		UserName:     info.Username,
		Introduction: "am a super administrator",
		RoleKeys:     roleKeys,
		Permissions:  permissions,
	}, nil
}

//获取个人信息
func (s *SysService) UserProfile(ctx context.Context, req *v1.UserInfoReq) (*v1.UserListData, error) {
	info, err := s.userc.GetUserProfile(ctx, req.UserId)
	if err != nil {
		return nil, errors.Unauthorized("获取用户信息失败", err.Error())
	}
	deptInfo, _ := s.deptc.GetInfoById(ctx, info.DeptId)
	if deptInfo == nil {
		deptInfo = &model.SysDept{}
	}
	out := &v1.UserListData{
		UserId:        info.UserId,
		Username:      info.Username,
		Mobile:        info.Mobile,
		Nickname:      info.Nickname,
		Status:        info.Status,
		Email:         info.Email,
		Address:       info.Address,
		Sex:           int32(info.Sex),
		Avatar:        info.Avatar,
		DeptId:        info.DeptId,
		IsAdmin:       info.IsAdmin,
		LastLoginIp:   info.LastLoginIp,
		LastLoginTime: global.GetDateByUnix(info.LastLoginTime),
		CreatedAt:     global.GetDateByUnix(info.CreatedAt),
		UpdatedAt:     global.GetDateByUnix(info.UpdatedAt),
		DeptInfo: &v1.UserListDataUserDept{
			DeptId:   deptInfo.DeptId,
			DeptName: deptInfo.DeptName,
		},
	}
	if len(info.Roles) > 0 {
		roles := make([]*v1.UserListDataUserRoles, len(info.Roles))
		for k, v := range info.Roles {
			roleInfo := &v1.UserListDataUserRoles{
				RoleId:  v.RoleId,
				Name:    v.Name,
				RoleKey: v.RoleKey,
			}
			roles[k] = roleInfo
		}
		out.Roles = roles
	}
	if len(info.Posts) > 0 {
		posts := make([]*v1.UserListDataUserPosts, len(info.Posts))
		for k, v := range info.Posts {
			postInfo := &v1.UserListDataUserPosts{
				PostId:   v.PostId,
				PostCode: v.PostCode,
				PostName: v.PostName,
			}
			posts[k] = postInfo
		}
		out.Posts = posts
	}
	return out, nil
}

func (s *SysService) UserList(ctx context.Context, req *v1.UserListReq) (*v1.UserListRes, error) {
	deptList, _, err := s.deptc.AllDeptList(ctx)
	if err != nil {
		return nil, errors.Unauthorized("获取部门列表失败", err.Error())
	}
	deptMap := make(map[int64]*model.SysDept)
	for _, v := range deptList {
		info := v
		deptMap[v.DeptId] = info
	}
	res, total, err := s.userc.ListUser(ctx, req)
	if err != nil {
		return nil, errors.Unauthorized("获取用户列表失败", err.Error())
	}
	list := make([]*v1.UserListData, len(res))
	for k, v := range res {
		userInfo := &v1.UserListData{
			UserId:        v.UserId,
			Username:      v.Username,
			Mobile:        v.Mobile,
			Nickname:      v.Nickname,
			Status:        v.Status,
			Email:         v.Email,
			Address:       v.Address,
			Sex:           int32(v.Sex),
			Avatar:        v.Avatar,
			Remark:        v.Remark,
			DeptId:        v.DeptId,
			IsAdmin:       v.IsAdmin,
			Birthday:      global.GetDateByUnix(v.Birthday),
			LastLoginIp:   v.LastLoginIp,
			LastLoginTime: global.GetDateByUnix(v.LastLoginTime),
			CreatedAt:     global.GetDateByUnix(v.CreatedAt),
			UpdatedAt:     global.GetDateByUnix(v.UpdatedAt),
		}
		if deptInfo, ok := deptMap[v.DeptId]; ok {
			log.Println("=====", deptInfo)
			userInfo.DeptInfo = &v1.UserListDataUserDept{
				DeptId:   deptInfo.DeptId,
				DeptName: deptInfo.DeptName,
			}
		}

		if v.Roles != nil {
			roleList := make([]*v1.UserListDataUserRoles, len(v.Roles))
			for key, val := range v.Roles {
				info := &v1.UserListDataUserRoles{
					RoleId:  val.RoleId,
					Name:    val.Name,
					RoleKey: val.RoleKey,
				}
				roleList[key] = info
			}
			userInfo.Roles = roleList
		}
		if v.Posts != nil {
			postList := make([]*v1.UserListDataUserPosts, len(v.Posts))
			for key, val := range v.Posts {
				info := &v1.UserListDataUserPosts{
					PostId:   val.PostId,
					PostName: val.PostName,
					PostCode: val.PostCode,
				}
				postList[key] = info
			}
			userInfo.Posts = postList

		}
		list[k] = userInfo
	}
	return &v1.UserListRes{Total: total, List: list}, nil
}

//更新用户状态
func (s *SysService) UpdateUserStatus(ctx context.Context, req *v1.UserStatusReq) (*v1.CommonRes, error) {
	err := s.userc.UpdateUserStatus(ctx, req.UserId, req.Status)
	return &v1.CommonRes{}, err
}

func (s *SysService) GetUserPostRoleList(ctx context.Context, _ *common.NullReq) (*v1.UserPostRoleListRes, error) {
	allRole, _ := s.rolec.AllListRole(ctx)
	allPost, _ := s.postc.AllListPost(ctx)
	allDept, _ := s.deptc.AllListDept(ctx)
	roleList := make([]*v1.UserPostRoleListResRoleInfo, len(allRole))
	postList := make([]*v1.UserPostRoleListResPostInfo, len(allPost))
	deptList := make([]*v1.UserPostRoleListResDeptInfo, len(allDept))
	if allRole != nil {
		for k, v := range allRole {
			info := &v1.UserPostRoleListResRoleInfo{}
			info.RoleId = v.RoleId
			info.Pid = v.PID
			info.Status = v.Status
			info.Sort = v.Sort
			info.RoleKey = v.RoleKey
			info.Name = v.Name
			roleList[k] = info
		}
	}
	if allPost != nil {
		for k, v := range allPost {
			info := &v1.UserPostRoleListResPostInfo{}
			info.PostId = v.PostId
			info.PostCode = v.PostCode
			info.PostName = v.PostName
			info.PostSort = v.PostSort
			postList[k] = info
		}
	}
	if allDept != nil {
		for k, v := range allDept {
			info := &v1.UserPostRoleListResDeptInfo{}
			info.DeptId = v.DeptId
			info.DeptName = v.DeptName
			info.ParentId = v.ParentId
			deptList[k] = info
		}
	}
	return &v1.UserPostRoleListRes{PostList: postList, RoleList: roleList, DeptList: deptList}, nil
}

func (s *SysService) UserDelete(ctx context.Context, req *v1.UserDeleteReq) (*v1.CommonRes, error) {
	err := s.userc.BatchDeleteUserByIds(ctx, req.UserIds)
	if err != nil {
		return nil, errors.Unauthorized("UserDeleteError", err.Error())
	}
	return &v1.CommonRes{}, nil
}

func (s *SysService) UserAdd(ctx context.Context, req *v1.UserAddReq) (*v1.CommonRes, error) {
	password := strings.TrimSpace(req.Password)
	insert := &model.SysUser{
		Username: req.Username,
		Nickname: req.Nickname,
		Mobile:   req.Mobile,
		Email:    req.Email,
		DeptId:   req.DeptId,
		Sex:      int(req.Sex),
		Avatar:   req.Avatar,
		Address:  req.Address,
		Remark:   req.Remark,
		IsAdmin:  int32(req.IsAdmin),
	}
	insert.CreateBy = req.CreateBy
	insert.CreatedAt = time.Now().Unix()
	insert.SetPassword(password)
	if len(req.PostIds) > 0 {
		insert.Posts, _ = s.postc.ListPostByIds(ctx, req.PostIds)
	}
	if len(req.RoleIds) > 0 {
		insert.Roles, _ = s.rolec.ListRoleByIds(ctx, req.RoleIds)
	}
	id, err := s.userc.LastIdCreateUser(ctx, insert)
	if err != nil {
		return nil, err
	}
	return &v1.CommonRes{Pong: fmt.Sprintf("%d", id)}, nil
}

func (s *SysService) UserUpdate(ctx context.Context, req *v1.UserUpdateReq) (*v1.CommonRes, error) {
	updateData := &model.SysUser{
		UserId:   req.UserId,
		Nickname: req.Nickname,
		Mobile:   req.Mobile,
		Email:    req.Email,
		DeptId:   req.DeptId,
		Sex:      int(req.Sex),
		Avatar:   req.Avatar,
		Address:  req.Address,
		Remark:   req.Remark,
		IsAdmin:  int32(req.IsAdmin),
	}
	updateData.UpdateBy = req.UpdateBy
	updateData.UpdatedAt = time.Now().Unix()
	if len(req.PostIds) > 0 {
		updateData.Posts, _ = s.postc.ListPostByIds(ctx, req.PostIds)
	}
	if len(req.RoleIds) > 0 {
		updateData.Roles, _ = s.rolec.ListRoleByIds(ctx, req.RoleIds)
	}
	err := s.userc.UpdateUser(ctx, updateData)
	if err != nil {
		return nil, err
	}
	return &v1.CommonRes{}, nil
}

func (s *SysService) LoginUserUpdate(ctx context.Context, req *v1.LoginUserUpdateReq) (*v1.CommonRes, error) {
	err := s.userc.LoginUpdateUser(ctx, req.UserId, req.LastLoginTime, req.LastLoginIp)
	if err != nil {
		return nil, errors.Unauthorized("LoginUserUpdateError", err.Error())
	}
	return &v1.CommonRes{}, nil
}

//修改用户密码
func (s *SysService) ChangePwd(ctx context.Context, req *v1.ChangePwdReq) (*v1.CommonRes, error) {
	userInfo, err := s.userc.GetInfoById(ctx, req.UserId)
	if err != nil {
		return nil, errors.Unauthorized("GetInfoByIdError", err.Error())
	}
	oldPwd := strings.TrimSpace(req.OldPassword)
	if !userInfo.CheckPassword(oldPwd) {
		return nil, errors.Unauthorized("CheckPasswordError", "老密码错误")
	}
	//设置新密钥密码 更新
	newPwd := strings.TrimSpace(req.NewPassword)
	userInfo.SetPassword(newPwd)
	if err := s.userc.ChangePwd(ctx, req.UserId, userInfo.Salt, userInfo.Password); err != nil {
		return nil, err
	}
	return &v1.CommonRes{}, nil
}

func (s *SysService) UserProfileSet(ctx context.Context, req *v1.UserProfileSetReq) (*v1.CommonRes, error) {
	data := make(map[string]interface{})
	data["mobile"] = req.Mobile
	data["nickname"] = req.Nickname
	data["email"] = req.Email
	data["avatar"] = req.Avatar
	data["sex"] = req.Sex
	data["update_by"] = req.UpdateBy
	err := s.userc.UpdateUserProfile(ctx, req.UserId, data)
	if err != nil {
		return nil, err
	}
	return &v1.CommonRes{}, nil
}
