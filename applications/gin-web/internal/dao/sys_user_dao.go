package dao

import (
	"context"

	"github.com/mathiasXie/curd"
	"github.com/mathiasXie/gin-web/applications/gin-web/internal/model"
	"github.com/mathiasXie/gin-web/pkg/logger"
	"gorm.io/gorm"
)

type SysUserDAO struct {
	sysUserModel *curd.Model[model.SysUser]
	sysRoleModel *curd.Model[model.SysRole]
}

type UserDomain struct {
	*model.SysUser
	Role *model.SysRole
}

func NewSysUserDAO(db *gorm.DB) *SysUserDAO {
	return &SysUserDAO{
		sysUserModel: curd.NewModel[model.SysUser](db),
		sysRoleModel: curd.NewModel[model.SysRole](db),
	}
}

func (u *SysUserDAO) AddUser(ctx context.Context, data map[string]interface{}) (int32, error) {
	insert, err := u.sysUserModel.Insert(ctx, data)
	if err != nil {
		return 0, err
	}
	return insert.Id, nil
}

func (u *SysUserDAO) GetAllUser(ctx context.Context) ([]*UserDomain, int64, error) {
	users, total, err := u.sysUserModel.SelectAll(ctx)
	if err != nil {
		return nil, 0, err
	}

	// 获取所有角色ID
	roleIds := make([]int64, len(users))
	for i, user := range users {
		roleIds[i] = int64(user.RoleId)
	}

	// 批量查询角色
	roles, _, err := u.sysRoleModel.Select(ctx, &curd.ModelBO{
		Conditions: []map[string]interface{}{
			{"id": roleIds},
		},
	})
	if err != nil {
		return nil, 0, err
	}

	// 将角色ID映射到角色对象
	roleMap := make(map[int64]*model.SysRole)
	for _, role := range roles {
		roleMap[role.Id] = role
	}

	resp := make([]*UserDomain, len(users))
	for i, user := range users {
		resp[i] = &UserDomain{
			SysUser: user,
			Role:    roleMap[int64(user.RoleId)],
		}
	}

	return resp, total, nil
}

func (u *SysUserDAO) GetUser(ctx context.Context, Id int64) (*model.SysUser, error) {
	user, err := u.sysUserModel.FindOne(ctx, Id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *SysUserDAO) UpdateUser(ctx context.Context, Id int64, data map[string]interface{}) error {
	err := u.sysUserModel.Update(ctx, Id, data)
	if err != nil {
		logger.CtxError(ctx, "[SysUserDAO]UpdateUser Error:", err.Error())
		return err
	}
	return nil
}

func (u *SysUserDAO) FindUser(ctx context.Context, userName string) (*model.SysUser, error) {
	user, total, err := u.sysUserModel.Select(ctx, &curd.ModelBO{
		Conditions: []map[string]interface{}{
			{"username": userName},
		},
	})
	if err != nil {
		logger.CtxError(ctx, "[SysUserDAO]FindUser Error:", err.Error())
		return nil, err
	}
	if total == 0 {
		return nil, nil
	}
	return user[0], nil
}

func (u *SysUserDAO) GetRole(ctx context.Context, Id int64) (*model.SysRole, error) {
	role, err := u.sysRoleModel.FindOne(ctx, Id)
	if err != nil {
		logger.CtxError(ctx, "[SysUserDAO]GetRole Error:", err.Error())
		return nil, err
	}
	return role, nil
}
