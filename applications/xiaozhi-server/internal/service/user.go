package service

import (
	"context"
	"errors"

	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/dto"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/internal/consts"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/internal/model"
	"github.com/mathiasXie/gin-web/pkg/logger"
	"gorm.io/gorm"
)

type UserService struct {
	ctx context.Context
	db  *gorm.DB
}

func NewUserService(ctx context.Context, db *gorm.DB) *UserService {
	return &UserService{ctx: ctx, db: db}
}

func (s *UserService) GetUserByRoleId(roleId int) (*dto.UserInfo, error) {
	var userRole model.AiRole
	if err := s.db.Where("id = ?", roleId).First(&userRole).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, consts.RetRoleNotFound
		}
		logger.CtxError(s.ctx, "GetUserByRoleId", err)
		return nil, err
	}

	var user model.AiUser
	if err := s.db.Where("id = ?", userRole.UserId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, consts.RetUserNotFound
		}
		logger.CtxError(s.ctx, "GetUserByRoleId", err)
		return nil, err
	}

	return &dto.UserInfo{
		ID:       int(user.Id),
		Username: user.UserName,
		Email:    user.Email,
		Role: &dto.UserRole{
			LLM:        userRole.Llm,
			LLMModelId: userRole.LlmModelId,
			TTS:        userRole.Tts,
			TTSVoiceId: userRole.TtsVoiceId,
			Language:   userRole.Language,
		},
	}, nil
}
