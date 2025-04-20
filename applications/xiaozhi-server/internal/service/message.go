package service

import (
	"context"

	"gorm.io/gorm"
)

type MessageService struct {
	ctx context.Context
	db  *gorm.DB
}

func NewMessageService(ctx context.Context, db *gorm.DB) *MessageService {
	return &MessageService{ctx: ctx, db: db}
}
