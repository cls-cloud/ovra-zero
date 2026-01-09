package online

import (
	"context"
	"errors"
	"fmt"
	"ovra/app/system/internal/logic/monitor/logininfor"
	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"
	"ovra/toolkit/auth"
	"ovra/toolkit/errx"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
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

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	l.ctx = ctx

	var (
		cursor  uint64
		allKeys []string
	)

	pattern := "token:*"

	for {
		keys, nextCursor, err := l.svcCtx.Rds.ScanCtx(ctx, cursor, pattern, 1000)
		if err != nil {
			return nil, err
		}
		allKeys = append(allKeys, keys...)
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	if len(allKeys) == 0 {
		resp.Rows = nil
		resp.Total = 0
		return resp, nil
	}

	type LoginInfo struct {
		Token string
		Cur   string
		Act   string
	}

	loginInfoMap := make(map[string]LoginInfo)

	cmds := make([]*redis.SliceCmd, 0, len(allKeys))

	err = l.svcCtx.Rds.PipelinedCtx(ctx, func(pipe redis.Pipeliner) error {
		for _, key := range allKeys {
			cmd := pipe.HMGet(
				ctx,
				key,
				auth.FieldLoginInfoId,
				auth.FieldToken,
				auth.FieldCurrentTime,
				auth.FieldActiveTimeout,
			)
			cmds = append(cmds, cmd)
		}
		return nil
	})
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	for _, cmd := range cmds {
		vals, err := cmd.Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			continue
		}

		if len(vals) != 4 {
			continue
		}

		infoID, _ := vals[0].(string)
		token, _ := vals[1].(string)
		curStr, _ := vals[2].(string)
		actStr, _ := vals[3].(string)

		if infoID == "" || token == "" {
			continue
		}

		if judgmentActive(curStr, actStr) {
			continue
		}

		loginInfoMap[infoID] = LoginInfo{
			Token: token,
			Cur:   curStr,
			Act:   actStr,
		}
	}

	if len(loginInfoMap) == 0 {
		resp.Rows = nil
		resp.Total = 0
		return resp, nil
	}

	loginInfoIds := make([]string, 0, len(loginInfoMap))
	for id := range loginInfoMap {
		loginInfoIds = append(loginInfoIds, id)
	}

	q := l.svcCtx.Dal.Query
	do := q.SysLogininfor.WithContext(ctx)

	if req.Ipaddr != "" {
		do = do.Where(q.SysLogininfor.Ipaddr.Like(fmt.Sprintf("%%%s%%", req.Ipaddr)))
	}
	if req.UserName != "" {
		do = do.Where(q.SysLogininfor.UserName.Like(fmt.Sprintf("%%%s%%", req.UserName)))
	}

	loginInfoList, err := do.
		Select(
			q.SysLogininfor.InfoID,
			q.SysLogininfor.UserName,
			q.SysLogininfor.Ipaddr,
			q.SysLogininfor.LoginLocation,
			q.SysLogininfor.LoginTime,
			q.SysLogininfor.Browser,
			q.SysLogininfor.Os,
			q.SysLogininfor.DeviceType,
			q.SysLogininfor.ClientKey,
		).
		Where(q.SysLogininfor.InfoID.In(loginInfoIds...)).
		Order(q.SysLogininfor.LoginTime.Desc()).
		Find()

	if err != nil {
		return nil, errx.GORMErr(err)
	}

	onlineList := make([]*types.OnlineInfo, 0, len(loginInfoList))

	for _, info := range loginInfoList {
		loginInfo, ok := loginInfoMap[info.InfoID]
		if !ok {
			continue
		}

		uc, err := auth.AnalyseToken(
			loginInfo.Token,
			l.svcCtx.Config.JwtAuth.AccessSecret,
		)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				continue
			}
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
			Token:         loginInfo.Token,
		})
	}

	resp.Rows = onlineList
	resp.Total = int64(len(onlineList))
	return resp, nil
}

// judgmentActive 判断是否活跃
func judgmentActive(curStr, actStr string) bool {

	curInt, _ := strconv.ParseInt(curStr, 10, 64)
	actInt, _ := strconv.ParseInt(actStr, 10, 64)
	now := time.Now().Unix()

	if curInt > 0 && actInt > 0 && now > curInt+actInt {
		return true // token 过期
	}
	return false
}
