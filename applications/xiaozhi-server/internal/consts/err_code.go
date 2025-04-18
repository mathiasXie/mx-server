package consts

type ErrCode struct {
	raw     error // raw err
	ErrCode string
	// 默认中文描述
	ErrMsg string
}

var (
	RetSuccess = &ErrCode{nil, "CW0000", "成功"}

	//	2000-3000 下游异常
	RetRpcError = &ErrCode{nil, "CW2000", "系统调用异常:%s"}

	// 3000-4000 中间件异常
	RetRedisError       = &ErrCode{nil, "CW3002", "redis连接异常"}
	RetAbaseError       = &ErrCode{nil, "CW3003", "abase连接异常"}
	RetMysqlError       = &ErrCode{nil, "CW3004", "mysql连接异常"}
	RetMysqlDupKeyError = &ErrCode{nil, "CC3005", "mysql唯一键冲突"}
)
