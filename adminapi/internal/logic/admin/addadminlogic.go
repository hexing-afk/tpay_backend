package admin

import (
	"context"
	"strings"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"
	"tpay_backend/model"

	"github.com/tal-tech/go-zero/core/logx"
)

type AddAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) AddAdminLogic {
	return AddAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddAdminLogic) AddAdmin(req types.AddAdminRequest) (*types.AddAdminResponse, error) {
	//解密
	plainPassword, err := common.DecryptPassword(req.Password)
	if err != nil {
		l.Errorf("密码解密发生错误,err:%v, password:%v", err, req.Password)
		return nil, common.NewCodeError(common.SysDBErr)
	}

	admin := &model.Admin{
		Username:     req.Username,
		Password:     common.CreateAdminPassword(plainPassword),
		EnableStatus: model.AdminEnableStatus,
		Phone:        req.Phone,
		Email:        req.Email,
	}

	err = model.NewAdminModel(l.svcCtx.DbEngine).Insert(admin)
	if err != nil {
		if strings.Contains(err.Error(), model.AdminUserUniqueErr.Error()) {
			l.Logger.Errorf("管理员[%v]已经存在", req.Username)
			return nil, common.NewCodeError(common.AccountRepeat)
		} else {
			l.Logger.Errorf("添加管理员[%v]失败, err=%v", req.Username, err)
			return nil, common.NewCodeError(common.SysDBAdd)
		}
	}

	return &types.AddAdminResponse{AdminId: admin.Id}, nil
}
