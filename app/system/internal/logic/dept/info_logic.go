package dept

import (
	"context"
	"encoding/json"
	"github.com/jinzhu/copier"
	"reflect"
	"system/internal/svc"
	"system/internal/types"
	"toolkit/errx"

	"github.com/zeromicro/go-zero/core/logx"
)

type InfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoLogic {
	return &InfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InfoLogic) Info(req *types.IdReq) (resp map[string]interface{}, err error) {
	res := new(types.DeptBase)
	q := l.svcCtx.Dal.Query
	sysPost, err := q.SysDept.WithContext(l.ctx).Where(q.SysDept.DeptID.Eq(req.Id)).First()
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	err = copier.Copy(&res, sysPost)
	if sysPost.Leader != "" {
		res.Leader = &sysPost.Leader
	} else {
		res.Leader = nil
	}
	// 特殊处理parentId 兼容前端
	var toMap map[string]interface{}
	bs, err := json.Marshal(res)
	if err != nil {
		return nil, errx.BizErr("json 序列化失败")
	}
	err = json.Unmarshal(bs, &toMap)
	if err != nil {
		return nil, errx.BizErr("json 反序列化失败")
	}
	if res.ParentID == "0" {
		toMap["parentId"] = 0
	}
	resp = toMap
	return
}

func structToMap(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	// 处理指针类型
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	data := make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// 跳过未导出字段
		if field.PkgPath != "" {
			continue
		}

		data[field.Name] = value.Interface()
	}
	return data
}
