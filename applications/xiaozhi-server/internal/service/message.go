package service

import (
	"context"
	"time"

	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/internal/model"
	"gorm.io/gorm"
)

type MessageService struct {
	ctx context.Context
	db  *gorm.DB
}

func NewMessageService(ctx context.Context, db *gorm.DB) *MessageService {
	return &MessageService{ctx: ctx, db: db}
}

func (s *MessageService) StoreChatRecord(userId int, DeviceId int, role string, text string) error {
	return s.db.Create(&model.AiMessage{
		UserId:    int32(userId),
		DeviceId:  int32(DeviceId),
		Role:      role,
		Message:   text,
		CreatedAt: time.Now(),
	}).Error
}

func (s *MessageService) GetChatRecords(userId int, deviceId int, limit int) ([]*model.AiMessage, error) {
	var records []*model.AiMessage
	err := s.db.Where("user_id = ? AND device_id = ?", userId, deviceId).Order("created_at DESC").Limit(limit).Find(&records).Error
	return records, err
}
