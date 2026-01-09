package monitor

import (
	"context"
	"errors"
	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type RedisMonitorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRedisMonitorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RedisMonitorLogic {
	return &RedisMonitorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RedisMonitorLogic) RedisMonitor() (resp *types.RedisMonitorResp, err error) {
	script := `return redis.call("INFO", "ALL")`

	evalResult, err := l.svcCtx.Rds.Eval(script, []string{})
	if err != nil {
		return nil, err
	}
	infoStr, ok := evalResult.(string)
	if !ok {
		return nil, errors.New("eval result not string")
	}

	resp, err = ParseRedisInfo(infoStr)
	if err != nil {
		return nil, err
	}
	return
}

func ParseRedisInfo(infoStr string) (*types.RedisMonitorResp, error) {
	resp := &types.RedisMonitorResp{
		Info:         make(map[string]string),
		CommandStats: make([]types.CommandState, 0),
	}

	lines := strings.Split(infoStr, "\n")

	var maxDbKeys int64 = 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "db") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				infoParts := strings.Split(parts[1], ",")
				for _, p := range infoParts {
					if strings.HasPrefix(p, "keys=") {
						keyStr := strings.TrimPrefix(p, "keys=")
						keysNum, err := strconv.ParseInt(keyStr, 10, 64)
						if err == nil && keysNum > maxDbKeys {
							maxDbKeys = keysNum
						}
					}
				}
			}
			continue
		}

		if strings.HasPrefix(line, "cmdstat_") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				continue
			}
			cmdName := strings.TrimPrefix(parts[0], "cmdstat_")
			fields := strings.Split(parts[1], ",")
			for _, f := range fields {
				if strings.HasPrefix(f, "calls=") {
					valStr := strings.TrimPrefix(f, "calls=")
					val, err := strconv.ParseInt(valStr, 10, 64)
					if err == nil {
						resp.CommandStats = append(resp.CommandStats, types.CommandState{
							Name:  cmdName,
							Value: strconv.FormatInt(val, 10),
						})
					}
					break
				}
			}
			continue
		}
		kv := strings.SplitN(line, ":", 2)
		if len(kv) == 2 {
			resp.Info[kv[0]] = kv[1]
		}
	}

	resp.DbSize = maxDbKeys

	return resp, nil
}
