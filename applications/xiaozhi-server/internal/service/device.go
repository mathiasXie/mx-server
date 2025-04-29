package service

import (
	"context"
	"errors"
	"math/rand"

	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/dto"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/internal/model"
	"github.com/mathiasXie/gin-web/pkg/logger"
	"gorm.io/gorm"
)

type DeviceService struct {
	ctx context.Context
	db  *gorm.DB
}

func NewDeviceService(ctx context.Context, db *gorm.DB) *DeviceService {
	return &DeviceService{ctx: ctx, db: db}
}

func (s *DeviceService) CreateDevice(device *model.AiDevice) error {
	return s.db.Create(device).Error
}

func (s *DeviceService) GetDeviceByMac(mac string) (*model.AiDevice, error) {
	var device model.AiDevice
	if err := s.db.Where("device_mac = ?", mac).First(&device).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.CtxError(s.ctx, "GetDeviceByMac", err)
		return nil, err
	}
	return &device, nil
}

func (s *DeviceService) GetDevice(device *model.AiDevice) error {
	return s.db.Where("id = ?", device.Id).First(device).Error
}

func (s *DeviceService) UpdateDevice(device *model.AiDevice) error {
	return s.db.Model(&model.AiDevice{}).Where("id = ?", device.Id).Updates(device).Error
}

func (s *DeviceService) GenerateBindCode(reqDevice *dto.OtaRequest, ip string) int {
	//先检查设备是否存在
	device, err := s.GetDeviceByMac(reqDevice.MacAddress)
	if err != nil {
		logger.CtxError(s.ctx, "[GenerateBindCode]GetDeviceByMac error", err)
		return 0
	}

	if device == nil {
		bindCode := s.checkBindCode()
		//创建设备
		device = &model.AiDevice{
			DeviceMac:     reqDevice.MacAddress,
			BindCode:      int32(bindCode),
			BoardSsid:     reqDevice.Board.SSID,
			BoardType:     reqDevice.Board.Type,
			BoardIp:       reqDevice.Board.IP,
			Language:      reqDevice.Language,
			Version:       reqDevice.Application.Version,
			ChipModelName: reqDevice.ChipModelName,
			Ip:            ip,
		}
		err = s.CreateDevice(device)
		if err != nil {
			logger.CtxError(s.ctx, "[GenerateBindCode]GenerateBindCode error", err)
			return 0
		}
	}
	return int(device.BindCode)
}

func (s *DeviceService) checkBindCode() int {

	code := rand.Intn(1000000)

	var device model.AiDevice
	if err := s.db.Where("bind_code = ?", code).First(&device).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return code
	}
	return s.checkBindCode()
}

func (s *DeviceService) GetDeviceList(device *model.AiDevice) ([]*model.AiDevice, error) {
	var devices []*model.AiDevice
	if err := s.db.Where("user_id = ?", device.UserId).Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

func (s *DeviceService) BindDevice(userId int, deviceId int, roleId int) error {

	err := s.UpdateDevice(&model.AiDevice{
		UserId:   int32(userId),
		BindCode: 0,
		Id:       int32(deviceId),
		RoleId:   int32(roleId),
	})
	if err != nil {
		logger.CtxError(s.ctx, "[BindDevice]BindDevice UpdateDevice error", err)
		return err
	}

	return nil
}

func (s *DeviceService) UnbindDevice(deviceId int) error {
	return s.UpdateDevice(&model.AiDevice{
		UserId:   0,
		RoleId:   0,
		BindCode: int32(s.checkBindCode()),
		Id:       int32(deviceId),
	})
}

func (s *DeviceService) GetDeviceByBindCode(bindCode int) (*model.AiDevice, error) {
	var device model.AiDevice
	if err := s.db.Where("bind_code = ?", bindCode).First(&device).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.CtxError(s.ctx, "[GetDeviceByBindCode]GetDeviceByBindCode error", err)
		return nil, err
	}
	return &device, nil
}
