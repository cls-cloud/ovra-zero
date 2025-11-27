package tenant

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"system/internal/dal/model"
	"system/internal/logic/dept"
	"system/internal/logic/user"
	"time"
	"toolkit/errx"
	"toolkit/utils"

	"system/internal/logic/role"
	"system/internal/svc"
	"system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddLogic {
	return &AddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddLogic) Add(req *types.ModifyTenantReq) error {
	q := l.svcCtx.Dal.Query
	//获取当前最大的租户编号
	tenantID, err := l.generateUniqueTenantID()
	if err != nil {
		return errx.BizErr("租户编号生成失败")
	}
	expireTime, err := time.Parse(time.DateTime, req.ExpireTime)
	if err != nil {
		return errx.BizErr("时间格式错误")
	}
	tenant := &model.SysTenant{
		ID:              utils.GetID(),
		TenantID:        tenantID,
		ContactUserName: req.ContactUserName,
		ContactPhone:    req.ContactPhone,
		CompanyName:     req.CompanyName,
		LicenseNumber:   req.LicenseNumber,
		Intro:           req.Intro,
		Domain:          req.Domain,
		Remark:          req.Remark,
		PackageID:       req.PackageId,
		ExpireTime:      expireTime,
		AccountCount:    req.AccountCount,
		Address:         req.Address,
	}
	if err := q.SysTenant.WithContext(l.ctx).Create(tenant); err != nil {
		return errx.GORMErr(err)
	}
	err = l.createAdminUser(tenant, req)
	if err != nil {
		return err
	}
	return nil
}

func (l *AddLogic) generateUniqueTenantID() (string, error) {
	q := l.svcCtx.Dal.Query
	rand.NewSource(time.Now().UnixNano())
	for {
		tenantId := fmt.Sprintf("%06d", rand.Intn(1000000))
		count, err := q.SysTenant.WithContext(l.ctx).Where(q.SysTenant.TenantID.Eq(tenantId)).Count()
		if err != nil {
			return "", errx.GORMErr(err)
		}
		if count == 0 {
			return tenantId, nil
		}
	}
}

func (l *AddLogic) createAdminUser(tenant *model.SysTenant, req *types.ModifyTenantReq) error {
	q := l.svcCtx.Dal.Query
	tenantPackage, err := q.SysTenantPackage.WithContext(l.ctx).Where(q.SysTenantPackage.PackageID.Eq(tenant.PackageID)).First()
	if err != nil {
		return errx.GORMErr(err)
	}
	//新增租户角色
	menuIds := strings.Split(tenantPackage.MenuIds, ",")
	roleId := utils.GetID()
	if err = role.NewAddLogic(l.ctx, l.svcCtx).Add(&types.AddOrUpdateRoleReq{
		RoleBase: types.RoleBase{
			RoleID:    roleId,
			RoleName:  "系统管理员",
			RoleKey:   "admin",
			RoleSort:  1,
			Status:    "0",
			Remark:    tenant.CompanyName,
			DataScope: "1",
			TenantID:  tenant.TenantID,
		},
		MenuIds: menuIds,
	}); err != nil {
		return err
	}
	//新增部门
	deptId := utils.GetID()
	err = dept.NewAddLogic(l.ctx, l.svcCtx).Add(&types.ModifyDeptReq{
		DeptBase: types.DeptBase{
			DeptID:   deptId,
			DeptName: tenant.CompanyName,
			Status:   "0",
			TenantID: tenant.TenantID,
			Phone:    tenant.ContactPhone,
			ParentID: "0",
			OrderNum: 1,
		},
	})
	//新增用户
	if err := user.NewAddUserLogic(l.ctx, l.svcCtx).AddUser(&types.AddOrUpdateUserReq{
		UserBase: types.UserBase{
			UserName: req.UserName,
			NickName: "管理员",
			Password: req.Password,
			Status:   "0",
			TenantID: tenant.TenantID,
			DeptID:   deptId,
		},
		RoleIds: []string{roleId},
	}); err != nil {
		return err
	}
	return nil
}
