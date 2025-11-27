package notice

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"system/internal/dal/model"
	"system/internal/svc"
	"system/internal/types"
	"time"
	"toolkit/auth"
	"toolkit/errx"

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

func (l *PageSetLogic) PageSet(req *types.PageSetNoticeReq) (resp *types.PageSetNoticeResp, err error) {
	offset := (req.PageNum - 1) * req.PageSize
	q := l.svcCtx.Dal.Query
	var result []struct {
		model.SysNotice
		CreateByName string `gorm:"column:user_name"`
	}
	do := q.SysNotice.WithContext(l.ctx).
		Select(q.SysNotice.ALL, q.SysUser.UserName).
		LeftJoin(q.SysUser, q.SysUser.UserID.EqCol(q.SysNotice.CreateBy))
	if req.NoticeTitle != "" {
		do = do.Where(q.SysNotice.NoticeTitle.Like(fmt.Sprintf("%%%s%%", req.NoticeTitle)))
	}
	if req.NoticeType != "" {
		do = do.Where(q.SysNotice.NoticeType.Eq(req.NoticeType))
	}
	if req.CreateBy != "" {
		do.Where(q.SysUser.UserName.Eq(fmt.Sprintf("%%%s%%", req.CreateBy)))
	}
	total, err := do.Count()
	tenantId := auth.GetTenantId(l.ctx)
	if tenantId != "" {
		do = do.Where(q.SysNotice.TenantID.Eq(tenantId))
	}
	err = do.Order(q.SysNotice.CreateTime.Desc()).Offset(int(offset)).Limit(int(req.PageSize)).Scan(&result)
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	resp = new(types.PageSetNoticeResp)
	resp.Total = total
	list := make([]*types.NoticeBase, len(result))
	for i, item := range result {
		list[i] = new(types.NoticeBase)
		if err = copier.Copy(&list[i], item); err != nil {
			return nil, err
		}
		list[i].CreateTime = item.CreateTime.Format(time.DateTime)
		list[i].NoticeContent = string(item.NoticeContent)
		list[i].CreateByName = item.CreateByName
	}
	resp.Rows = list
	return
}
