package merchant

import (
	"context"
	"fmt"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"
	"tpay_backend/utils"

	"gorm.io/gorm"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ModifyMerchantBalanceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewModifyMerchantBalanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) ModifyMerchantBalanceLogic {
	return ModifyMerchantBalanceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ModifyMerchantBalanceLogic) ModifyMerchantBalance(adminId int64, req types.ModifyMerchantBalanceRequest) error {
	merchant, err := model.NewMerchantModel(l.svcCtx.DbEngine).FindOneById(req.MerchantId)
	if err != nil {
		if err == model.ErrRecordNotFound {
			l.Errorf("商户[%v]不存在, err=%v", req.MerchantId, err)
			return common.NewCodeError(common.UserNotExist)
		} else {
			l.Errorf("查询商户[%v]失败, err=%v", req.MerchantId, err)
			return common.NewCodeError(common.SysDBUpdate)
		}
	}

	var description = ""

	switch req.OpType {
	case model.OpTypeAddBalance:
		description = fmt.Sprintf("管理员[%v]给商户[%v]增加了[%v]余额", adminId, req.MerchantId, req.ChangeAmount)

		txErr := l.svcCtx.DbEngine.Transaction(func(tx *gorm.DB) error {
			// 余额变动的额外信息
			walletLogExt := model.WalletLogExt{
				BusinessNo: "",                         // 业务单号
				Source:     model.AmountSourcePlatform, // 变动来源：1-手动调账
				Remark:     req.Remark,                 // 备注
			}
			// 1.增加商户金额
			if err := model.NewMerchantModel(tx).PlusBalance(merchant.Id, req.ChangeAmount, walletLogExt); err != nil {
				return err
			}

			return nil
		})
		if txErr != nil {
			l.Errorf("增加商户[%v]余额失败, err=%v", req.MerchantId, err)
			return common.NewCodeError(common.SysDBUpdate)
		}
	case model.OpTypeMinusBalance:
		description = fmt.Sprintf("管理员[%v]给商户[%v]减少了[%v]余额", adminId, req.MerchantId, req.ChangeAmount)

		if req.ChangeAmount > merchant.Balance {
			l.Errorf("商户[%v]余额不足, balance=%v", req.MerchantId, merchant.Balance)
			return common.NewCodeError(common.InsufficientBalance)
		}

		txErr := l.svcCtx.DbEngine.Transaction(func(tx *gorm.DB) error {
			// 余额变动的额外信息
			walletLogExt := model.WalletLogExt{
				BusinessNo: "",                         // 业务单号
				Source:     model.AmountSourcePlatform, // 变动来源：1-手动调账
				Remark:     req.Remark,                 // 备注
			}

			// 1.减商户金额
			if err := model.NewMerchantModel(tx).MinusBalance(merchant.Id, req.ChangeAmount, walletLogExt); err != nil {
				return err
			}

			return nil
		})
		if txErr != nil {
			l.Errorf("减少商户[%v]余额失败, err=%v", req.MerchantId, err)
			return common.NewCodeError(common.SysDBUpdate)
		}

	default:
		l.Errorf("req.OpType[%v]参数有误", req.OpType)
		return common.NewCodeError(common.InvalidParam)
	}

	log := &model.AdminWebLog{
		LogNo:       utils.GetDailyId(),
		AdminId:     adminId,
		Description: description,
		Type:        model.LogTypeMerchant,
	}
	if err = model.NewAdminWebLogModel(l.svcCtx.DbEngine).Insert(log); err != nil {
		l.Errorf("记录管理员[%v]操作日志失败, err=%v", adminId, err)
	}

	return nil
}
