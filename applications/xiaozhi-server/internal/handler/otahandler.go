package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/dto"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/internal/model"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/internal/service"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/loader"
	"github.com/mathiasXie/gin-web/config"
	"github.com/mathiasXie/gin-web/pkg/logger"
	"github.com/mathiasXie/gin-web/utils"
)

func OtaHandler(ctx *gin.Context) {

	deviceId := ctx.Request.Header.Get("Device-Id")
	if deviceId == "" {
		ctx.JSON(http.StatusOK, gin.H{"error": "Device ID is required"})
		return
	}
	var deviceInfo dto.OtaRequest
	if err := ctx.ShouldBindJSON(&deviceInfo); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	otaResponse := dto.OtaResponse{}

	// 生成websocket地址
	scheme := "ws"
	defaultPort := 80
	if config.Instance.Server.Ssl {
		scheme = "wss"
		defaultPort = 443
	}
	var url string
	if config.Instance.Server.Port == defaultPort {
		url = fmt.Sprintf("%s://%s", scheme, config.Instance.Server.Host)
	} else {
		url = fmt.Sprintf("%s://%s:%d", scheme, config.Instance.Server.Host, config.Instance.Server.Port)
	}
	url = fmt.Sprintf("%s/xiaozhi/v1/", url)
	otaResponse.WebSocket = dto.WebSocket{
		URL: url,
	}

	// 判断客户端版本是否有升级 TODO
	if deviceInfo.Application.Version == "" {
		otaResponse.Firmware = dto.Firmware{
			Version: "", //可升级的版本号
			URL:     "", // 固件下载URL
		}
	} else {
		otaResponse.Firmware = dto.Firmware{
			Version: deviceInfo.Application.Version,
			URL:     "",
		}
	}

	// 服务器时间
	_, offsetSeconds := time.Now().Zone()
	otaResponse.ServerTime = dto.ServerTime{
		Timestamp:      time.Now().UnixMilli(),
		TimezoneOffset: int32(offsetSeconds / 60),
	}

	// 查询用户激活情况
	deviceService := service.NewDeviceService(ctx, loader.GetDB(ctx, true))
	device, err := deviceService.GetDeviceByMac(deviceInfo.MacAddress)
	if err != nil {
		logger.CtxError(ctx, "OTA[GetDeviceByMac] error", err.Error())
		ctx.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	if device == nil || device.RoleId == 0 {
		bindCode := deviceService.GenerateBindCode(&deviceInfo, utils.GetClientIP(ctx))
		otaResponse.Activation = &dto.Activation{
			Code:      strconv.Itoa(bindCode),
			Message:   fmt.Sprintf("%s\n%d", "激活验证码", bindCode),
			Challenge: deviceInfo.MacAddress,
		}
	} else {
		go deviceService.UpdateDevice(&model.AiDevice{
			DeviceMac:     deviceInfo.MacAddress,
			BoardSsid:     deviceInfo.Board.SSID,
			BoardType:     deviceInfo.Board.Type,
			BoardIp:       deviceInfo.Board.IP,
			Language:      deviceInfo.Language,
			Version:       deviceInfo.Application.Version,
			ChipModelName: deviceInfo.ChipModelName,
			Ip:            utils.GetClientIP(ctx),
			Id:            device.Id,
		})
	}
	ctx.JSON(http.StatusOK, otaResponse)
}

func ActivateHandler(ctx *gin.Context) {
	deviceId := ctx.Request.Header.Get("Device-Id")
	if deviceId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Device ID is required"})
		return
	}

	clientId := ctx.Request.Header.Get("Client-Id")
	if clientId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Client ID is required"})
		return
	}

	deviceService := service.NewDeviceService(ctx, loader.GetDB(ctx, true))
	device, err := deviceService.GetDeviceByMac(deviceId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if device == nil || device.RoleId == 0 {
		ctx.JSON(http.StatusAccepted, gin.H{"message": "Device activation timeout"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Device activated", "device_id": device.Id})
}
