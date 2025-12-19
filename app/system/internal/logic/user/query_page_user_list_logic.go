package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"ovra/app/system/internal/dal/model"
	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"
	"ovra/toolkit/auth"
	"ovra/toolkit/errx"
	"time"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type QueryPageUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryPageUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryPageUserListLogic {
	return &QueryPageUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type ListResp struct {
	Total int64                    `json:"total"`
	Rows  []map[string]interface{} `json:"rows"`
}

func (l *QueryPageUserListLogic) QueryPageUserList(req *types.QueryPageUserListReq) (resp *ListResp, err error) {
	sysUser := l.svcCtx.Dal.Query.SysUser
	sysDept := l.svcCtx.Dal.Query.SysDept

	do := sysUser.WithContext(l.ctx).
		LeftJoin(sysDept, sysDept.DeptID.EqCol(sysUser.DeptID)).
		Select(
			sysUser.ALL,
			sysDept.DeptName, // 添加 dept_name 字段
		)

	// 条件查询
	if req.DeptId != "" {
		deptIds, err := GetDeptIds(req, l.svcCtx.Db, l.ctx)
		if err != nil {
			return nil, errx.GORMErr(err)
		}
		do = do.Where(sysUser.DeptID.In(deptIds...))
	}
	if req.UserName != "" {
		do = do.Where(sysUser.UserName.Like(fmt.Sprintf("%%%s%%", req.UserName)))
	}
	if req.NickName != "" {
		do = do.Where(sysUser.NickName.Like(fmt.Sprintf("%%%s%%", req.NickName)))
	}
	if req.PhoneNumber != "" {
		do = do.Where(sysUser.Phonenumber.Like(fmt.Sprintf("%%%s%%", req.PhoneNumber)))
	}
	if req.Status != "" {
		do = do.Where(sysUser.Status.Eq(req.Status))
	}
	if req.BeginTime != "" {
		beginTime, err := time.Parse(time.DateTime, req.BeginTime)
		if err != nil {
			return nil, errors.New("invalid beginTime format, expected: YYYY-MM-DD HH:mm:ss")
		}
		do = do.Where(sysUser.CreateTime.Gte(beginTime))
	}
	if req.EndTime != "" {
		endTime, err := time.Parse(time.DateTime, req.EndTime)
		if err != nil {
			return nil, errors.New("invalid endTime format, expected: YYYY-MM-DD HH:mm:ss")
		}
		do = do.Where(sysUser.CreateTime.Lte(endTime))
	}

	// 分页
	total, err := do.Count()
	if err != nil {
		return nil, err
	}

	offset := (req.PageNum - 1) * req.PageSize
	var result []struct {
		model.SysUser
		DeptName string `gorm:"column:dept_name"`
	}
	tenantId := auth.GetTenantId(l.ctx)
	if tenantId != "" {
		do = do.Where(sysUser.TenantID.Eq(tenantId))
	}
	err = do.Offset(int(offset)).Limit(int(req.PageSize)).Scan(&result)
	if err != nil {
		return nil, err
	}

	// 特殊处理userId 兼容前端
	resp = new(ListResp)
	resp.Total = total
	resp.Rows = make([]map[string]interface{}, len(result))
	for i, row := range result {
		var userDetail types.UserAndDeptName
		_ = copier.Copy(&userDetail, row.SysUser)
		userDetail.CreateTime = row.CreateTime.Format(time.DateTime)
		userDetail.DeptName = row.DeptName
		userDetail.UserID = row.UserID
		userDetail.DeptID = row.DeptID
		var toMap map[string]interface{}
		bs, err := json.Marshal(userDetail)
		if err != nil {
			return nil, errx.BizErr("json 序列化失败")
		}
		err = json.Unmarshal(bs, &toMap)
		if err != nil {
			return nil, errx.BizErr("json 反序列化失败")
		}
		resp.Rows[i] = toMap
		if userDetail.UserID == "1" {
			resp.Rows[i]["userId"] = 1
		}
	}

	return
}

func GetDeptIds(req *types.QueryPageUserListReq, db *gorm.DB, ctx context.Context) ([]string, error) {
	var deptIds []string
	err := db.WithContext(ctx).
		Table("sys_dept").
		Select("dept_id").
		Where("FIND_IN_SET(?, ancestors)", req.DeptId).
		Find(&deptIds).Error
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	deptIds = append(deptIds, req.DeptId)
	return deptIds, nil
}
