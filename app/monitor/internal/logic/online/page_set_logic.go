package online

import (
	"context"
	"fmt"
	"ovra/app/monitor/internal/logic/logininfor"
	"ovra/app/monitor/internal/svc"
	"ovra/app/monitor/internal/types"
	"ovra/toolkit/auth"
	"ovra/toolkit/errx"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type PageSetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPageSetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageSetLogic {
	return &PageSetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
func (l *PageSetLogic) PageSet(req *types.OnlineQuery) (resp *types.PageSetOnlineResp, err error) {
	resp = new(types.PageSetOnlineResp)
	//evalResult, err := l.svcCtx.Rds.EvalCtx(l.ctx, "return redis.call('KEYS', KEYS[1])", []string{"token:*"})
	//if err != nil {
	//	return nil, err
	//}
	//
	//// 类型断言
	//results, ok := evalResult.([]interface{})
	//if !ok {
	//	return nil, errors.New("eval result not []interface{}")
	//}
	//
	//// 解析 key 列表
	//var keys []string
	//for _, v := range results {
	//	if s, ok := v.(string); ok {
	//		keys = append(keys, s)
	//	}
	//}
	var cursor uint64
	var allKeys []string

	for {
		keys, nextCursor, err := l.svcCtx.Rds.ScanCtx(l.ctx, cursor, "token:*", 100)
		if err != nil {
			return nil, err
		}
		allKeys = append(allKeys, keys...)
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	type TokenInfo struct {
		LoginInfoId string
		Token       string
	}
	loginInfoMap := make(map[string]string)

	for _, key := range allKeys {
		infoId, err := l.svcCtx.Rds.HgetCtx(l.ctx, key, auth.FieldLoginInfoId)
		if err != nil {
			if err.Error() == "redis: nil" {
				continue
			}
			return nil, err
		}

		token, err := l.svcCtx.Rds.HgetCtx(l.ctx, key, auth.FieldToken)
		if err != nil {
			if err.Error() == "redis: nil" {
				continue
			}
			return nil, err
		}

		loginInfoMap[infoId] = token
	}

	// 提取所有 loginInfoId
	var loginInfoIds []string
	for id := range loginInfoMap {
		loginInfoIds = append(loginInfoIds, id)
	}
	q := l.svcCtx.Dal.Query
	do := q.SysLogininfor.WithContext(l.ctx)
	if req.Ipaddr != "" {
		do = do.Where(q.SysLogininfor.Ipaddr.Like(fmt.Sprintf("%%%s%%", req.Ipaddr)))
	}
	if req.UserName != "" {
		do = do.Where(q.SysLogininfor.UserName.Like(fmt.Sprintf("%%%s%%", req.UserName)))
	}
	loginInfoList, err := do.
		Where(q.SysLogininfor.InfoID.In(loginInfoIds...)).
		Order(q.SysLogininfor.LoginTime.Desc()).
		Find()
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	resp.Total = int64(len(loginInfoList))
	onlineList := make([]*types.OnlineInfo, 0, len(loginInfoList))

	for _, info := range loginInfoList {
		token := loginInfoMap[info.InfoID]
		uc, err := auth.AnalyseToken(token, l.svcCtx.Config.JwtAuth.AccessSecret)
		if err != nil {
			return nil, errx.BizErr("系统异常")
		}
		onlineList = append(onlineList, &types.OnlineInfo{
			Browser:       info.Browser,
			ClientKey:     info.ClientKey,
			DeviceType:    info.DeviceType,
			DeptName:      uc.DeptName,
			Ipaddr:        logininfor.DesensitizeIP(info.Ipaddr),
			LoginLocation: logininfor.DesensitizeLocation(info.LoginLocation),
			LoginTime:     info.LoginTime.Format(time.DateTime),
			Os:            info.Os,
			UserName:      info.UserName,
			Token:         token,
		})
	}
	resp.Rows = onlineList
	return
}
