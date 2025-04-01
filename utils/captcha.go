package utils

import (
	"context"
	"fmt"
	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"
	"image/color"
	"time"
)

func GenerateCaptcha(ctx context.Context, redis *redis.Client) (string, string, error) {
	//定义一个driver
	var driver base64Captcha.Driver
	//创建一个字符串类型的验证码驱动DriverString, DriverChinese :中文驱动
	driverString := base64Captcha.DriverString{
		Height:          80,                                 //高度
		Width:           200,                                //宽度
		NoiseCount:      0,                                  //干扰数
		ShowLineOptions: 4,                                  //展示个数
		Length:          4,                                  //长度
		Source:          "23456789qwertyuplkjhgfdsazxcvbnm", //验证码随机字符串来源
		BgColor: &color.RGBA{ // 背景颜色
			R: 3,
			G: 102,
			B: 100,
			A: 125,
		},
		Fonts: []string{"chromohv.ttf"}, // 字体
	}
	driver = driverString.ConvertFonts()
	//var store = base64Captcha.DefaultMemStore

	//配置RedisStore, RedisStore实现base64Captcha.Store接口
	var store base64Captcha.Store = &RedisStore{redisClient: redis, ctx: ctx}
	//生成验证码
	c := base64Captcha.NewCaptcha(driver, store)
	id, b64s, _, err := c.Generate()
	return id, b64s, err
}

// VerifyCaptcha 校验验证码
func VerifyCaptcha(id string, VerifyValue string, redis *redis.Client) bool {
	var store = &RedisStore{redisClient: redis}
	// 参数说明: id 验证码id, verifyValue 验证码的值, true: 验证成功后是否删除原来的验证码
	if store.Verify(id, VerifyValue, true) {
		return true
	} else {
		return false
	}
}

const (
	prefixCaptcha = "CaptchaCode"
	expireCaptcha = time.Hour
)

type RedisStore struct {
	redisClient *redis.Client
	ctx         context.Context
}

func (s *RedisStore) Set(id, value string) error {
	key := fmt.Sprintf("%s:%s", prefixCaptcha, id)
	return s.redisClient.Set(s.ctx, key, value, expireCaptcha).Err()
}

func (s *RedisStore) Get(id string, clear bool) string {

	key := fmt.Sprintf("%s:%s", prefixCaptcha, id)
	val := s.redisClient.Get(context.Background(), key).Val()

	if clear {
		_ = s.redisClient.Del(context.Background(), key)
	}
	return val
}

func (s *RedisStore) Verify(id, answer string, clear bool) bool {
	return s.Get(id, clear) == answer
}
