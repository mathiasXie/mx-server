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

// 查找用户的主要方法
func (s *UserService) GetUserByPhoneEmailOrUsername(username string) (*model.AiUser, error) {
	var user model.AiUser
	var err error

	// 尝试通过手机号码查找用户
	if user, err = s.findUserByPhone(username); err == nil {
		return &user, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.CtxError(s.ctx, "[UserService]GetUserByPhoneEmailOrUsername phone", err)
		return nil, err
	}

	// 尝试通过邮箱查找用户
	if user, err = s.findUserByEmail(username); err == nil {
		return &user, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.CtxError(s.ctx, "[UserService]GetUserByPhoneEmailOrUsername email", err)
		return nil, err
	}

	// 尝试通过用户名查找用户
	if user, err = s.findUserByUsername(username); err == nil {
		return &user, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.CtxError(s.ctx, "[UserService]GetUserByPhoneEmailOrUsername username", err)
		return nil, err
	}

	return nil, consts.RetUserNotFound
}

// 通过手机号码查找用户
func (s *UserService) findUserByPhone(phone string) (model.AiUser, error) {
	var user model.AiUser
	err := s.db.Where("phone = ?", phone).First(&user).Error
	return user, err
}

// 通过邮箱查找用户
func (s *UserService) findUserByEmail(email string) (model.AiUser, error) {
	var user model.AiUser
	err := s.db.Where("email = ?", email).First(&user).Error
	return user, err
}

// 通过用户名查找用户
func (s *UserService) findUserByUsername(username string) (model.AiUser, error) {
	var user model.AiUser
	err := s.db.Where("user_name = ?", username).First(&user).Error
	return user, err
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

func (s *UserService) CreateUser(user model.AiUser) error {
	return s.db.Create(&user).Error
}

func (s *UserService) CheckUserExists(user dto.SignUpRequest) string {

	if _, err := s.findUserByPhone(user.Phone); err == nil {
		return "phone"
	}
	if _, err := s.findUserByEmail(user.Email); err == nil {
		return "email"
	}
	if _, err := s.findUserByUsername(user.Username); err == nil {
		return "username"
	}
	return ""
}
