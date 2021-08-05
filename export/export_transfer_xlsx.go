package export

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/tealeg/xlsx"
	"strconv"
	"time"
)

type CreateTransferOrderFileRequest struct {
	Sheet           string              // 工作表
	Timezone        string              // 时区
	Title           string              // 文件标题
	Total           *TransferOrderTotal // 合计数据
	Content         []TransferOrder     // 订单数据
	IsDivideHundred bool                // 金额是否需要除以100，币种单位为分的都需要
}

type TransferOrderTotal struct {
	Total               int64
	TotalReqAmount      int64
	TotalFee            int64
	TotalIncreaseAmount int64
}

type TransferOrder struct {
	Id              int64
	MerchantName    string
	OrderNo         string
	MerchantOrderNo string
	ReqAmount       int64
	Rate            float64
	SingleFee       int64
	Fee             int64
	IncreaseAmount  int64
	ChannelName     string
	OrderStatus     string
	CreateTime      int64
	UpdateTime      int64
}

var transferOrderFileHeadName = []string{
	"商户名称",
	"平台订单号",
	"商户订单号",
	"订单金额",
	"手续费",
	"费率",
	"到账金额",
	"通道",
	"状态",
	"创建时间",
	"更新时间",
}

// 创建代付订单Excel文件
func CreateTransferOrderFile(file *xlsx.File, req *CreateTransferOrderFileRequest) (*xlsx.File, error) {
	var (
		orderLis = req.Content
		total    = req.Total
	)

	// 创建文件
	sheet, err := file.AddSheet(req.Sheet)
	if err != nil {
		return nil, err
	}

	// 设置标题样式
	titleStyle := xlsx.NewStyle()
	titleStyle.Font = xlsx.Font{
		Size:      18,
		Name:      "宋体",
		Bold:      true,
		Italic:    false,
		Underline: false,
	}
	titleStyle.Alignment = xlsx.Alignment{
		Horizontal: "center",
		Vertical:   "center",
	}

	// 添加标题
	titleRow := sheet.AddRow()
	titleRow.SetHeightCM(2)
	titleCell := titleRow.AddCell()
	titleCell.HMerge = 20
	titleCell.Value = fmt.Sprintf("%s(时区：%s)", req.Title, req.Timezone)
	titleCell.SetStyle(titleStyle)

	// 添加文件头
	headStyle := xlsx.NewStyle()
	headStyle.Font = xlsx.Font{
		Size: 13,
		Name: "宋体",
		Bold: true,
	}

	headRow := sheet.AddRow()
	headRow.SetHeightCM(0.7)
	for _, v := range transferOrderFileHeadName {
		cell := headRow.AddCell()
		cell.SetStyle(headStyle)
		cell.Value = v
	}

	// 写入数据
	for _, order := range orderLis {
		contentRow := sheet.AddRow()
		contentRow.SetHeightCM(0.7)

		// 商户名称
		merchantName := contentRow.AddCell()
		merchantName.Value = order.MerchantName

		// 平台订单号
		orderNo := contentRow.AddCell()
		orderNo.Value = order.OrderNo

		// 商户订单号
		mchOrderNo := contentRow.AddCell()
		mchOrderNo.Value = order.MerchantOrderNo

		// 订单请求金额
		reqAmount := contentRow.AddCell()

		// 手续费
		fee := contentRow.AddCell()

		// 费率
		rate := contentRow.AddCell()

		// 到账金额
		increaseAmount := contentRow.AddCell()

		var singleFee string
		if req.IsDivideHundred {
			reqAmount.Value = decimal.NewFromInt(order.ReqAmount).Div(decimal.NewFromInt(100)).Round(2).String()

			fee.Value = decimal.NewFromInt(order.Fee).Div(decimal.NewFromInt(100)).Round(2).String()

			increaseAmount.Value = decimal.NewFromInt(order.IncreaseAmount).Div(decimal.NewFromInt(100)).Round(2).String()

			if order.SingleFee > 0 {
				singleFee = decimal.NewFromInt(order.SingleFee).Div(decimal.NewFromInt(100)).Round(2).String()
			}
		} else {
			reqAmount.SetInt64(order.ReqAmount)

			fee.SetInt64(order.Fee)

			increaseAmount.SetInt64(order.IncreaseAmount)

			if order.SingleFee > 0 {
				singleFee = strconv.FormatInt(order.SingleFee, 10)
			}
		}

		if singleFee != "" {
			rate.Value = fmt.Sprintf("%v", order.Rate) + "% + " + singleFee
		} else {
			rate.Value = fmt.Sprintf("%v", order.Rate) + "%"
		}

		// 通道
		channel := contentRow.AddCell()
		channel.Value = order.ChannelName

		// 状态
		status := contentRow.AddCell()
		status.Value = order.OrderStatus

		// 创建时间
		createTime := contentRow.AddCell()
		createTime.Value = time.Unix(order.CreateTime, 0).Format(TimeFormat)

		// 更新时间
		updateTime := contentRow.AddCell()
		updateTime.Value = time.Unix(order.UpdateTime, 0).Format(TimeFormat)
	}

	sheet.AddRow().SetHeightCM(0.7)

	// 合计
	countStyle := xlsx.NewStyle()
	countStyle.Font = xlsx.Font{
		Size: 13,
		Name: "宋体",
		Bold: true,
	}
	countStyle.Alignment = xlsx.Alignment{
		Horizontal: "center",
		Vertical:   "center",
	}

	countRow := sheet.AddRow()
	countRow.SetHeightCM(0.7)
	countCell := countRow.AddCell()
	countCell.HMerge = 1
	countCell.SetStyle(countStyle)
	countCell.Value = "合计"
	countRow.AddCell()
	countRow.AddCell()

	if req.IsDivideHundred {
		countRow.AddCell().Value = decimal.NewFromInt(total.TotalReqAmount).Div(decimal.NewFromInt(100)).Round(2).String()

		countRow.AddCell().Value = decimal.NewFromInt(total.TotalFee).Div(decimal.NewFromInt(100)).Round(2).String()
		countRow.AddCell()

		countRow.AddCell().Value = decimal.NewFromInt(total.TotalIncreaseAmount).Div(decimal.NewFromInt(100)).Round(2).String()
	} else {
		countRow.AddCell().Value = strconv.FormatInt(total.TotalReqAmount, 10)

		countRow.AddCell().Value = strconv.FormatInt(total.TotalFee, 10)
		countRow.AddCell()

		countRow.AddCell().Value = strconv.FormatInt(total.TotalIncreaseAmount, 10)
	}

	return file, nil
}
