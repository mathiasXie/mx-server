package service

import (
	"context"
	"errors"

	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/internal/model"
	"github.com/mathiasXie/gin-web/pkg/logger"
	"gorm.io/gorm"
)

type RoleService struct {
	ctx context.Context
	db  *gorm.DB
}

func NewRoleService(ctx context.Context, db *gorm.DB) *RoleService {
	return &RoleService{
		ctx: ctx,
		db:  db,
	}
}

func (s *RoleService) GetRoleByID(id int) (*model.AiRole, error) {
	var role model.AiRole
	if err := s.db.Where("id = ?", id).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.CtxError(s.ctx, "GetRoleByID", err)
		return nil, err
	}
	return &role, nil
}

func (s *RoleService) GetRolesByUserId(userId int) ([]*model.AiRole, error) {
	var roles []*model.AiRole
	if err := s.db.Where("user_id = ?", userId).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}
