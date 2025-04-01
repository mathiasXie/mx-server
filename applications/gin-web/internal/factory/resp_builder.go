package factory

import (
	"context"
	"reflect"

	"github.com/mathiasXie/gin-web/applications/gin-web/internal/consts"
)

type CommonResp interface {
	SetRetCode(val string)
	SetRetMsg(val string)
}

// InitResp 生成非空resp
func InitResp[T any](resp T) T {
	v := reflect.ValueOf(resp)
	if v.Kind() != reflect.Ptr || !v.IsNil() {
		return resp
	}
	v = reflect.New(v.Type().Elem())
	return v.Interface().(T)
}

// BuildCommonResp 构建返回值公共字段（需保证resp非空）
func BuildCommonResp[T CommonResp](ctx context.Context, resp T, err error) T {
	resp = InitResp(resp)
	if err == nil {
		resp.SetRetCode(consts.SuccessCode)
		resp.SetRetMsg("成功")

		return resp
	}

	//logs.CtxError(ctx, "Got err: %+v", err)

	//resp.SetRetCode(string(werr.GetCode(err)))
	//resp.SetRetMsg(werr.GetMessage(err))

	return resp
}
