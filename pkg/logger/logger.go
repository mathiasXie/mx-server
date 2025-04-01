package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mathiasXie/gin-web/config"
	"github.com/mathiasXie/gin-web/consts"
	"github.com/mathiasXie/gin-web/pkg/logger/log"
	"github.com/mathiasXie/gin-web/utils"
)

var (
	logHandler log.Logger
	logLocker  sync.Mutex
)

func InitLog() {
	logLocker.Lock()
	if logHandler == nil {
		logHandler = log.NewLogrusLogger(config.Instance)
	}
	logLocker.Unlock()
}

func CtxInfo(ctx context.Context, data ...interface{}) {
	logHandler.Info(ctx, data)
}

func CtxError(ctx context.Context, data ...interface{}) {
	logHandler.Error(ctx, data)
}

func CtxPanic(ctx context.Context, data ...interface{}) {
	logHandler.Panic(ctx, data)
}
func CtxWarn(ctx context.Context, data ...interface{}) {
	logHandler.Warn(ctx, data)
}

func CtxDebug(ctx context.Context, data ...interface{}) {
	logHandler.Debug(ctx, data)
}

type ResponseLogger struct {
	gin.ResponseWriter
	bodyBytes []byte
}

func (r *ResponseLogger) Write(b []byte) (int, error) {
	r.bodyBytes = append(r.bodyBytes, b...)
	return r.ResponseWriter.Write(b)
}

func HttpLog(conf config.LoggerConf) gin.HandlerFunc {
	return func(c *gin.Context) {
		var logId = c.GetHeader(consts.LogID)
		if logId == "" {
			logId = utils.GenerateLogId()
			c.Request.Header.Add(consts.LogID, logId)
		}
		c.Set(consts.LogID, logId)
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path

		if conf.LogIDShowHeader {
			c.Header(consts.LogID, logId)
		}

		// 获取请求参数
		var requestData interface{}
		if c.Request.Method != "GET" {
			// 对于非GET请求，读取body
			bodyBytes, _ := c.GetRawData()
			if len(bodyBytes) > 0 {
				// 尝试解析为JSON
				var jsonData interface{}
				if err := json.Unmarshal(bodyBytes, &jsonData); err == nil {
					requestData = utils.MaskSensitiveData(jsonData)
				} else {
					requestData = string(bodyBytes)
				}
				// 重新设置body，因为GetRawData会清空body
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		} else {
			// 对于GET请求，获取URL参数并脱敏
			queryParams := c.Request.URL.Query()
			maskedQuery := make(map[string]interface{})
			for key, values := range queryParams {
				if len(values) == 1 {
					maskedQuery[key] = utils.MaskSensitiveData(values[0])
				} else {
					maskedQuery[key] = utils.MaskSensitiveData(values)
				}
			}
			requestData = maskedQuery
		}

		// 脱敏请求头
		maskedHeaders := make(map[string][]string)
		for key, values := range c.Request.Header {
			if strings.ToLower(key) == "authorization" {
				maskedHeaders[key] = []string{"******"}
			} else {
				maskedHeaders[key] = values
			}
		}

		if !slices.Contains(conf.SkipPaths, path) {
			CtxInfo(c, "HttpEnter", map[string]interface{}{
				"Path":        path,
				"Method":      c.Request.Method,
				"ClientIP":    utils.GetClientIP(c),
				"RequestData": requestData,
				"Headers":     maskedHeaders,
			})
		}

		responseLogger := &ResponseLogger{ResponseWriter: c.Writer}
		c.Writer = responseLogger

		// Process request
		c.Next()

		if !slices.Contains(conf.SkipPaths, path) {
			// 脱敏响应数据
			var responseData interface{}
			if len(responseLogger.bodyBytes) > 0 {
				if err := json.Unmarshal(responseLogger.bodyBytes, &responseData); err == nil {
					responseData = utils.MaskSensitiveData(responseData)
				} else {
					responseData = string(responseLogger.bodyBytes)
				}
			}

			CtxInfo(c, "HttpOut", map[string]interface{}{
				"Latency":      time.Since(start).Milliseconds(),
				"BodySize":     c.Writer.Size(),
				"StatusCode":   c.Writer.Status(),
				"ResponseData": responseData,
			})
		}
	}
}
