package export

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/tealeg/xlsx"
	"strconv"
	"time"
)

type CreateWalletLogFileRequest struct {
	Sheet           string      // 工作表
	Title           string      // 文件标题
	Timezone        string      // 时区
	IsDivideHundred bool        // 金额是否需要除以100，币种单位为分的都需要
	Content         []WalletLog // 订单数据
}

type WalletLog struct {
	Id           int64
	BusinessNo   string
	ChangeAmount int64
	AfterBalance int64
	OpType       string
	OrderType    string
	Remark       string
	CreateTime   int64
}

//钱包流水表头
var walletLogHead = []string{
	"业务单号",
	"变动金额",
	"变动后余额",
	"出入账类型",
	"订单类型",
	"备注",
	"创建时间",
}

func CreateWalletLogFile(file *xlsx.File, req *CreateWalletLogFileRequest) (*xlsx.File, error) {
	orderLis := req.Content

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
	for _, v := range walletLogHead {
		cell := headRow.AddCell()
		cell.SetStyle(headStyle)
		cell.Value = v
	}

	// 写入数据
	for _, order := range orderLis {
		contentRow := sheet.AddRow()
		contentRow.SetHeightCM(0.7)

		// 业务单号
		orderNo := contentRow.AddCell()
		orderNo.Value = order.BusinessNo

		// 变动金额
		changeAmount := contentRow.AddCell()

		// 变动后余额
		afterBalance := contentRow.AddCell()

		if req.IsDivideHundred {
			changeAmount.Value = decimal.NewFromInt(order.ChangeAmount).Div(decimal.NewFromInt(100)).Round(2).String()

			afterBalance.Value = decimal.NewFromInt(order.AfterBalance).Div(decimal.NewFromInt(100)).Round(2).String()
		} else {
			changeAmount.Value = strconv.FormatInt(order.ChangeAmount, 10)

			afterBalance.Value = strconv.FormatInt(order.AfterBalance, 10)
		}

		// 出入账类型
		opType := contentRow.AddCell()
		opType.Value = order.OpType

		// 订单类型
		orderType := contentRow.AddCell()
		orderType.Value = order.OrderType

		// 备注
		remark := contentRow.AddCell()
		remark.Value = order.Remark

		// 创建时间
		createTime := contentRow.AddCell()
		createTime.Value = time.Unix(order.CreateTime, 0).Format(TimeFormat)
	}

	return file, nil
}
