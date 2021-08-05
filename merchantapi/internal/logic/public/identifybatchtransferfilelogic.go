package public

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/tealeg/xlsx"
	"strings"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/model"
	"tpay_backend/utils"

	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type IdentifyBatchTransferFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIdentifyBatchTransferFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) IdentifyBatchTransferFileLogic {
	return IdentifyBatchTransferFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IdentifyBatchTransferFileLogic) IdentifyBatchTransferFile(merchantId int64, req types.IdentifyBatchTransferFileRequest) (*types.IdentifyBatchTransferFileReply, error) {

	//1.确认文件是否存在
	fileInfo, err := model.NewUploadFileLogModel(l.svcCtx.DbEngine).FindOne(req.FileName)
	if err != nil {
		l.Errorf("查询文件信息失败，err:%v", err)
		return nil, common.NewCodeError(common.FileNotExist)
	}

	if fileInfo.AccountId != merchantId || fileInfo.AccountType != model.UploadFileLogAccountTypeMerchant {
		l.Errorf("文件不是该商户的, fileInfo:%+v ; err:%v ", fileInfo, err)
		return nil, common.NewCodeError(common.FileNotExist)
	}

	//2.获取数据，分析、处理数据
	fileContentList, err := l.getBatchAnalysisData(req.FileName, merchantId)
	if err != nil {
		l.Errorf("处理数据发生错误, err:%v", err)
		return nil, err
	}

	//3.获取统计分析结果
	resultData, err := l.getBatchAnalysisResult(fileContentList)
	if err != nil {
		l.Errorf("处理数据发生错误, err:%v", err)
		return nil, err
	}

	resultData.MerchantId = merchantId

	//4.插入批量订单
	if err := model.NewTransferBatchOrderModel(l.svcCtx.DbEngine).Insert(resultData); err != nil {
		l.Errorf("插入批量订单失败,err:%v", err)
		return nil, common.NewCodeError(common.SysDBAdd)
	}

	return &types.IdentifyBatchTransferFileReply{
		BatchNo: resultData.BatchNo,
	}, nil
}

//	获取数据，分析、处理数据
func (l *IdentifyBatchTransferFileLogic) getBatchAnalysisData(fileName string, merchantId int64) (datas []*model.FileContent, errs error) {
	l.Infof("数据分析开始 start ")

	//从云存储获取文件
	bytes, err := l.svcCtx.CloudStorage.GetObject(fileName)
	if err != nil {
		l.Errorf("从云存储获取文件失败,FileName:%s, err:%v", fileName, err)
		return nil, common.NewCodeError(common.FileNotExist)
	}

	//分析文件信息，得到结果
	xlFile, err := xlsx.OpenBinary(bytes)
	if err != nil {
		l.Errorf("open failed: %s\n", err)
		return nil, common.NewCodeError(common.SysUnKnow)
	}

	//查询商户信息
	merchant, err := model.NewMerchantModel(l.svcCtx.DbEngine).FindOneById(merchantId)
	if err != nil {
		l.Errorf("查询商户信息出错, err:%v", err)
		return nil, common.NewCodeError(common.GetLoginUserInfoFailed)
	}

	//查询商家转账配置
	currencyInfo, err := model.NewCurrencyModel(l.svcCtx.DbEngine).FindByCurrency(merchant.Currency)
	if err != nil {
		l.Errorf("查询商户币种出错，err：%v", err)
		return nil, common.NewCodeError(common.SysDBErr)
	}

	var errArgs []string

	for _, sheet := range xlFile.Sheets {
		for j, row := range sheet.Rows {
			if j == 0 { //去掉标题
				continue
			}
			if len(row.Cells) == 0 { //去掉空行
				continue
			}
			data := &model.FileContent{
				Row:            fmt.Sprintf("%d", j), //行号
				ChannelCode:    "",                   //平台代付通道code
				AccountName:    "",                   //收款人姓名
				CardNumber:     "",                   //收款卡号
				BankName:       "",                   //银行名称
				BankBranchName: "",                   //支行名称
				Amount:         0,                    //金额
				Remark:         "",                   //备注
			}

			//序号	*收款人姓名	*收款卡号	*银行名称	支行名称	  *代付通道代码	*金额	备注
			for k, cell := range row.Cells {
				text := cell.String()
				switch k {
				case 0: //序号（这是商家批量转账的序号，后台不做处理）
				case 1: //收款人姓名
					data.AccountName = text
				case 2: //收款卡号
					data.CardNumber = text
				case 3: //银行名称
					data.BankName = text
				case 4: //支行名称
					data.BankBranchName = text
				case 5: //代付通道代码
					data.ChannelCode = text
				case 6: //金额
					amount, err := decimal.NewFromString(text)
					if err != nil {
						l.Errorf("行号%v, 金额错误,err:%v", j, err)
						continue
						//return nil, err
					}

					//2.35-->235
					if currencyInfo.IsDivideHundred == model.DivideHundred {
						//2.35 * 100
						amount = amount.Mul(decimal.NewFromInt(100))
					}

					//235
					data.Amount = amount.IntPart()

				case 7: //备注
					data.Remark = text
				}

				l.Infof("行号%v 第%v列, 内容: %s\n", j, k, text)
			}

			// 判断必填字段是否都必填了
			switch "" {
			case data.AccountName:
				fallthrough
			case data.CardNumber:
				fallthrough
			case data.BankName:
				fallthrough
			case data.ChannelCode:
				errArgs = append(errArgs, fmt.Sprintf("第%v行", j))
			default:
			}

			if data.Amount == 0 {
				errArgs = append(errArgs, fmt.Sprintf("第%v行", j))
			}

			datas = append(datas, data)
		}
	}

	if len(errArgs) != 0 {
		errTempStr := "文件%v有误，请修正后提交"
		return nil, errors.New(fmt.Sprintf(errTempStr, strings.Join(errArgs, "、")))
	}

	//
	l.Infof("数据分析完 end ")

	return datas, nil
}

//获取分析结果
func (l *IdentifyBatchTransferFileLogic) getBatchAnalysisResult(fileContentList []*model.FileContent) (*model.TransferBatchOrder, error) {
	resultData := &model.TransferBatchOrder{
		BatchNo:     utils.GetUniqueId(),   // 批量号
		TotalNumber: 0,                     // 订单总笔数
		TotalAmount: 0,                     // 订单总金额
		Status:      model.BatchStatusInit, // 批量状态
		//MerchantId:  "",                    // 商户id
		FileContent: "",                      // 文件内容json
		GenerateAll: model.GenerateAllUndone, // 是否已全部生成订单 1-完成，2-未完成
	}
	jsonByte, err := json.Marshal(fileContentList)
	if err != nil {
		l.Errorf("json编码文件内容失败, err:%v", err)
		return nil, err
	}
	resultData.FileContent = string(jsonByte)

	for _, data := range fileContentList {
		resultData.TotalNumber += 1
		resultData.TotalAmount += data.Amount
	}

	return resultData, nil
}
