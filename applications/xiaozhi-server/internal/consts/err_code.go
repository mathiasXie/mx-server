package consts

import "fmt"

type ErrCode struct {
	raw     error // raw err
	ErrCode string
	// 默认中文描述
	ErrMsg string
}

func (e *ErrCode) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrCode, e.ErrMsg)
}

var (
	RetSuccess = &ErrCode{nil, "CW0000", "成功"}

	RetRoleNotFound = &ErrCode{nil, "CW0001", "角色不存在"}

	RetUserNotFound = &ErrCode{nil, "CW0002", "用户不存在"}

	RetDeviceNotFound = &ErrCode{nil, "CW0003", "设备不存在"}

	RetDeviceBindCodeNotFound = &ErrCode{nil, "CW0004", "绑定码不存在"}

	RetTTSConfigError = &ErrCode{nil, "CW0005", "TTS配置错误"}

	RetLLMConfigError = &ErrCode{nil, "CW0006", "大模型配置错误"}
	//	2000-3000 下游异常
	RetRpcError = &ErrCode{nil, "CW2000", "系统调用异常:%s"}

	// 3000-4000 中间件异常
	RetRedisError       = &ErrCode{nil, "CW3002", "redis连接异常"}
	RetAbaseError       = &ErrCode{nil, "CW3003", "abase连接异常"}
	RetMysqlError       = &ErrCode{nil, "CW3004", "mysql连接异常"}
	RetMysqlDupKeyError = &ErrCode{nil, "CC3005", "mysql唯一键冲突"}
)
