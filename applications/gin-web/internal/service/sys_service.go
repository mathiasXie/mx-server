package service

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mathiasXie/gin-web/applications/gin-web/dto"
	"github.com/mathiasXie/gin-web/applications/gin-web/internal/dao"
	"github.com/mathiasXie/gin-web/applications/gin-web/internal/model"
	"github.com/mathiasXie/gin-web/applications/gin-web/loader"
	"github.com/mathiasXie/gin-web/pkg/logger"
	"github.com/mathiasXie/gin-web/utils"
)

type SysService struct {
	sysUserDAO  *dao.SysUserDAO
	sysRouteDAO *dao.SysRouteDAO
}

func NewSysService(ctx context.Context) *SysService {
	return &SysService{
		sysUserDAO:  dao.NewSysUserDAO(loader.GetDB(ctx)),
		sysRouteDAO: dao.NewSysRouteDAO(loader.GetDB(ctx)),
	}
}

func (s *SysService) Login(ctx *gin.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.sysUserDAO.FindUser(ctx, req.Username)
	if err != nil {
		logger.CtxError(ctx, "[SysService]Login error:", err)
		return nil, err
	}
	if user == nil {
		return nil, nil
	}
	if !utils.ComparePassword(user.Password, req.Password) {
		return nil, nil
	}
	err = s.sysUserDAO.UpdateUser(ctx, int64(user.Id), map[string]interface{}{
		"last_time": time.Now(),
		"last_ip":   utils.GetClientIP(ctx),
	})
	if err != nil {
		return nil, nil
	}
	token, _ := utils.MakeToken(uint64(user.Id), user.Username)
	return &dto.LoginResponse{
		UserId: user.Id,
		Token:  token,
	}, nil
}

func (s *SysService) FindOneUser(ctx *gin.Context, userId int64) (*model.SysUser, *model.SysRole, []*dao.UserRoleRoute, error) {
	user, err := s.sysUserDAO.GetUser(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "[SysService]FindOneUser error:", err)
		return nil, nil, nil, err
	}
	if user == nil {
		return nil, nil, nil, nil
	}
	role, _ := s.sysUserDAO.GetRole(ctx, int64(user.RoleId))

	roleRoute, _ := s.sysRouteDAO.GetRoleRoutes(ctx, int64(user.RoleId))

	return user, role, roleRoute, nil
}

func (s *SysService) GetUser(ctx *gin.Context) ([]*dto.UserListResponse, error) {
	users, _, err := s.sysUserDAO.GetAllUser(ctx)
	if err != nil {
		logger.CtxError(ctx, "[SysService]GetUser error:", err)
		return nil, err
	}
	if users == nil {
		return nil, nil
	}
	resp := make([]*dto.UserListResponse, len(users))
	for _, user := range users {
		resp = append(resp, &dto.UserListResponse{
			Id:       int64(user.Id),
			Username: user.Username,
			Realname: user.Realname,
			Email:    user.Email,
			Mobile:   user.Mobile,
			RoleID:   int64(user.RoleId),
			RoleName: user.Role.Name,
			Status:   user.Status,
			IsSuper:  user.IsSuper,
			LastTime: user.LastTime.Format(time.DateTime),
			LastIp:   user.LastIp,
		})
	}
	return resp, nil
}

func (s *SysService) UpdateUser(ctx *gin.Context, req *dto.UserUpdateRequest) error {
	data := map[string]interface{}{
		"username": req.Username,
		"realname": req.Realname,
		"email":    req.Email,
		"mobile":   req.Mobile,
		"role_id":  req.RoleId,
		"status":   req.Status,
		"is_super": req.IsSuper,
	}
	var err error
	if req.Id != 0 {
		err = s.sysUserDAO.UpdateUser(ctx, req.Id, data)
	} else {
		_, err = s.sysUserDAO.AddUser(ctx, data)
	}
	if err != nil {
		logger.CtxError(ctx, "[SysService]UpdateUser error:", err)
		return err
	}
	return nil
}

func (s *SysService) GetAllSysRoute(ctx *gin.Context) ([]*model.SysRoute, error) {
	routes, err := s.sysRouteDAO.GetSysRoute(ctx)
	if err != nil {
		logger.CtxError(ctx, "[SysService]GetSysRoute error:", err)
		return nil, err
	}
	return routes, nil
}

func (s *SysService) BuildMenuRoute(ctx *gin.Context, userRouteIds []int64, parentId int64) ([]*model.SysRoute, error) {
	routes, err := s.sysRouteDAO.BuildMenuRoute(ctx, userRouteIds, parentId)
	if err != nil {
		logger.CtxError(ctx, "[SysService]GetSysRoute error:", err)
		return nil, err
	}
	return routes, nil
}

func (s *SysService) GetRouteByParentID(ctx *gin.Context, parentId int64) ([]*model.SysRoute, error) {
	routes, err := s.sysRouteDAO.GetRouteByParentID(ctx, parentId)
	if err != nil {
		logger.CtxError(ctx, "[SysService]GetRouteByParentID error:", err)
		return nil, err
	}
	return routes, nil
}

func (s *SysService) RouteUpdate(ctx *gin.Context, req dto.RouteUpdateRequest) error {
	sysRoute := map[string]interface{}{
		"name":      req.Name,
		"label":     req.Label,
		"parent_id": req.ParentId,
		"path":      req.Path,
		"api_path":  req.ApiPath,
		"icon":      req.Icon,
		"sequence":  req.Sequence,
		"type":      req.Type,
		"status":    req.Status,
		"component": req.Component,
	}
	var err error
	if req.Id > 0 {
		err = s.sysRouteDAO.Update(ctx, req.Id, sysRoute)
	} else {
		err = s.sysRouteDAO.Insert(ctx, sysRoute)
	}
	return err
}

func (s *SysService) RouteDelete(ctx *gin.Context, req dto.RouteDeleteRequest) error {
	return s.sysRouteDAO.Delete(ctx, req.Id)
}
