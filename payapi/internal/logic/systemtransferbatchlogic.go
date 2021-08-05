package logic

import (
	"context"
	"fmt"
	"strings"
	"tpay_backend/model"
	"tpay_backend/payapi/internal/common"
	"tpay_backend/utils"

	"tpay_backend/payapi/internal/svc"
	"tpay_backend/payapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type SystemTransferBatchLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	merchant *model.Merchant
}

func NewSystemTransferBatchLogic(ctx context.Context, svcCtx *svc.ServiceContext, merchant *model.Merchant) SystemTransferBatchLogic {
	return SystemTransferBatchLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		merchant: merchant,
	}
}

func (l *SystemTransferBatchLogic) SystemTransferBatch(req types.SystemTransferBatchReq) (*types.SystemTransferBatchReply, error) {
	// 0.参数检查
	if strings.TrimSpace(req.BatchNo) == "" {
		return nil, common.NewCodeErrorWithMsg(common.VerifyParamFailed, "batch_no不能为空")
	}

	if l.merchant.Balance <= 0 {
		return nil, common.NewCodeError(common.InsufficientBalance)
	}

	// 1.查询批量付款订单
	batchOrder, err := model.NewTransferBatchOrderModel(l.svcCtx.DbEngine).FindByBatchNo(l.merchant.Id, req.BatchNo)
	if err != nil {
		if err == model.ErrRecordNotFound {
			l.Errorf("商户[%v]批量付款[%v]不存在", l.merchant.Id, req.BatchNo)
			return nil, common.NewCodeError(common.OrderNotExist)
		}
		l.Errorf("查询商户[%v]批量付款[%v]失败, err=%v", l.merchant.Id, req.BatchNo, err)
		return nil, common.NewCodeError(common.SystemInternalErr)
	}

	if batchOrder.Status == model.BatchStatusSuccess && batchOrder.GenerateAll == model.GenerateAllDone {
		l.Errorf("查询商户[%v]批量付款[%v]已完成", l.merchant.Id, req.BatchNo)
		return nil, common.NewCodeError(common.SystemInternalErr)
	}

	// 统计已生成订单，根据结果修改批量付款订单状态
	defer l.syncStatus(batchOrder)

	// 2.取出订单中的行号
	var rowNos []string
	fileMap := make(map[string]model.FileContent)
	for _, v := range batchOrder.FileContentSlice {
		rowNos = append(rowNos, v.Row)
		if _, ok := fileMap[v.Row]; !ok {
			fileMap[v.Row] = v
		} else {
			l.Errorf("查询商户[%v]批量付款[%v]有重复的文件数据", l.merchant.Id, req.BatchNo)
			l.Errorf("文件数据：%+v", v)
		}
	}

	// 3.从代付订单表中匹配，找出批量付款中未生成订单的行号
	notOrderRow, err := l.pickedNotExistOrder(l.merchant.MerchantNo, batchOrder.BatchNo, rowNos)
	if err != nil {
		return nil, common.NewCodeError(common.SystemInternalErr)
	}

	// 如果没有未生成订单的行号，表示批量付款已完成
	if len(notOrderRow) <= 0 {
		l.Errorf("查询商户[%v]批量付款[%v]已完成", l.merchant.Id, req.BatchNo)
		return &types.SystemTransferBatchReply{}, nil
	}

	// 4.根据行号拿到付款信息下单
	for _, row := range notOrderRow {
		transfer := NewTransferPlaceOrder(l.ctx, l.svcCtx, l.merchant)
		_, err := transfer.TransferPlaceOrder(TransferPlaceOrderRequest{
			Mode:               model.TransferModePro,
			MchOrderNo:         fmt.Sprintf("%s%s", "MBT", utils.GetDailyId()),
			Amount:             fileMap[row].Amount,
			Currency:           l.merchant.Currency,
			OrderSource:        model.TransferOrderSourceMerchantPayment,
			TradeType:          fileMap[row].ChannelCode,
			BankName:           fileMap[row].BankName,
			BankCardNo:         fileMap[row].CardNumber,
			BankCardHolderName: fileMap[row].AccountName,
			BankBranchName:     fileMap[row].BankBranchName,
			BankCode:           "",
			NotifyUrl:          "",
			ReturnUrl:          "",
			Attach:             "",
			Remark:             fileMap[row].Remark,
			BatchNo:            batchOrder.BatchNo,
			BatchRowNo:         row,
		})

		if err != nil {
			l.Errorf("批量付款[%v]批次的[%v]行数据下单失败: %v", req.BatchNo, row, err)
			return nil, common.NewCodeErrorWithMsg(common.OrderFailed, fmt.Sprintf("行号%v下单发生错误[%v]", row, err))
		}
	}

	return &types.SystemTransferBatchReply{}, nil
}

func (l *SystemTransferBatchLogic) pickedNotExistOrder(merchantNo string, batchNo string, rowNos []string) ([]string, error) {
	// 根据批次号和文件内容中的行号匹配是否已经生成了订单
	rows, err := model.NewTransferOrderModel(l.svcCtx.DbEngine).FindByBatchNoAndRowNo(merchantNo, batchNo, rowNos)
	if err != nil {
		logx.Errorf("查询代付订单失败, err=%v", err)
		return nil, err
	}

	var notOrderRow []string
	for _, v := range rowNos {
		if utils.InSlice(v, rows) {
			continue
		}

		notOrderRow = append(notOrderRow, v)
	}

	return notOrderRow, nil
}

func (l *SystemTransferBatchLogic) syncStatus(data *model.TransferBatchOrder) {
	// 4.订单生成完后判断该批次的已生成订单数量和总数量是否相等
	hadNum, err := model.NewTransferOrderModel(l.svcCtx.DbEngine).CountByBatchNo(l.merchant.MerchantNo, data.BatchNo)
	if err != nil {
		l.Errorf("查询批次订单[%v]已有订单数量失败, err=%v", data.Id, err)
		return
	}

	// 相等就修改transfer_batch_order表中的generate_all为1
	if hadNum == data.TotalNumber {
		err := model.NewTransferBatchOrderModel(l.svcCtx.DbEngine).UpdateStatus(data.Id, model.BatchStatusSuccess, model.GenerateAllDone)
		if err != nil {
			l.Errorf("修改批次订单Id[%v]状态失败, err=%v", data.Id, err)
			return
		}
		l.Infof("批次[%v]已完成代付订单生成, num:%v", data.BatchNo, hadNum)
		return
	}
}
