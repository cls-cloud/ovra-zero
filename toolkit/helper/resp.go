package helper

import (
	"context"
	"encoding/json"
	"net/http"
	"ovra/toolkit/errx"
	"reflect"

	"github.com/zeromicro/go-zero/core/logx"
)

type Response struct {
	Code int32       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func Success(data interface{}) any {
	if data == nil {
		return Response{
			Code: http.StatusOK,
			Msg:  "操作成功",
		}
	}

	// 判断是否有 Rows 字段
	val := reflect.Indirect(reflect.ValueOf(data))
	if val.Kind() == reflect.Struct {
		if val.FieldByName("Rows").IsValid() || val.FieldByName("Data").IsValid() {
			base := map[string]interface{}{
				"code": http.StatusOK,
				"msg":  "操作成功",
			}
			for k, v := range structToMap(data) {
				base[k] = v
			}
			return base
		}
	}

	return Response{
		Code: http.StatusOK,
		Msg:  "操作成功",
		Data: data,
	}
}

func Fail(err error) *Response {
	se := errx.FromError(err)
	return &Response{
		Code: se.Code,
		Msg:  se.Message,
		Data: nil,
	}
}

func OkHandler(_ context.Context, v interface{}) any {
	return Success(v)
}

func ErrHandler(name string) func(ctx context.Context, err error) (int, any) {
	return func(ctx context.Context, err error) (int, any) {
		logx.WithContext(ctx).Errorf("【%s】 err %v", name, err)
		return http.StatusOK, Fail(err)
	}
}

func structToMap(obj interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	if obj == nil {
		return m
	}
	b, _ := json.Marshal(obj)
	_ = json.Unmarshal(b, &m)
	return m
}
