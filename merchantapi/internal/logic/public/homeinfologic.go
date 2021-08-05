package public

import (
	"context"
	"time"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"
	"tpay_backend/model"

	"github.com/tal-tech/go-zero/core/logx"
)

type HomeInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	userId int64
}

func NewHomeInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext, userId int64) HomeInfoLogic {
	return HomeInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		userId: userId,
	}
}

func (l *HomeInfoLogic) HomeInfo() (*types.HomeInfoResponse, error) {
	merchant, err := model.NewMerchantModel(l.svcCtx.DbEngine).FindOneById(l.userId)
	if err != nil {
		l.Errorf("查询商户id出错, userId[%v], err[%v]", l.userId, err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	//当天
	toDay := time.Now().Format("2006-01-02")
	toDayCount, err := model.NewPayOrderModel(l.svcCtx.DbEngine).FindMerchantToDayCountByDay(merchant.MerchantNo)
	if err != nil {
		l.Errorf("查询当天[%v]统计数据出错, userId[%v], err[%v]", toDay, l.userId, err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	var receiveData []types.ReceiveData
	//
	for i := 7; i > 0; i-- {
		dayTime := time.Now().Add(time.Hour * -24 * time.Duration(i))
		successAmount, err := model.NewPayOrderModel(l.svcCtx.DbEngine).FindMerchantCountByDay(merchant.MerchantNo, dayTime.Format("2006-01-02"))
		if err != nil {
			l.Errorf("查询当天[%v]统计数据出错, userId[%v], err[%v]", dayTime.Format("2006-01-02"), l.userId, err)
			return nil, common.NewCodeError(common.SysDBGet)
		}

		receiveData = append(receiveData, types.ReceiveData{
			CreateTime: dayTime.Unix(),
			Amount:     successAmount,
		})
	}

	return &types.HomeInfoResponse{
		OrderNumber:        toDayCount.OrderNumber,        //今日收款总订单数
		SuccessOrderNumber: toDayCount.SuccessOrderNumber, //今日收款成功订单数
		SuccessAmount:      toDayCount.SuccessAmount,      //今日成功收款金额
		ReceiveList:        receiveData,                   //收款统计数据
	}, nil
}
