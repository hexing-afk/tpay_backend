package utils

import (
	"github.com/shopspring/decimal"
)

// 代收-计算商户手续费-默认内扣
func CalculatePayOrderFeeMerchant(reqAmount int64, singleFee int64, rate float64) int64 {
	total := decimal.NewFromInt(reqAmount)

	// 总金额 * (费率/100) + 单笔费用
	fee := total.Mul(decimal.NewFromFloat(rate)).Div(decimal.NewFromInt(100)).Add(decimal.NewFromInt(singleFee))

	// 四舍五入去除小数点
	fee = fee.Round(0)

	// 取整数
	return fee.IntPart()
}

// 代收-计算上游手续费-默认内扣
func CalculatePayOrderFeeUpstream(reqAmount int64, singleFee int64, rate float64) int64 {
	total := decimal.NewFromInt(reqAmount)

	// 总金额 * (费率/100) + 单笔费用
	fee := total.Mul(decimal.NewFromFloat(rate)).Div(decimal.NewFromInt(100)).Add(decimal.NewFromInt(singleFee))

	// 四舍五入去除小数点
	fee = fee.Round(0)

	// 取整数
	return fee.IntPart()
}

// 代付-手续费内扣方式-计算请求上游的金额
//
// payee_real_amount 收款方实际到账金额
// singleFee 上游单笔手续费
// rate 上游的手续费率
func CalculateUpstreamInnerAmount(payeeRealAmount int64, singleFee int64, rate float64) int64 {
	// 请求上游的金额 = 收款方实际到账金额 + (请求上游的金额 * 上游的手续费率 + 上游单笔手续费)

	// 示例：
	// 请求上游的金额: A
	// 收款方实际到账金额: 95
	// 上游的手续费率: 3%
	// 上游单笔手续费: 1
	// 计算公式: 95 + (A*3% + 1) = A
	// 计算公式: 95 + 1 + A*3% = A
	// 计算公式: 96 + A*3% = A
	// 计算公式: 96 = A - A*3%
	// 计算公式: 96 = A(1 - 3%)
	// 计算公式: 96 = A(97%)
	// 计算公式: 96 = A * 97 / 100
	// 计算公式: 96 * 100 / 97 = A
	// 计算公式: 96 * 100 / (100-3) = A
	amount := decimal.NewFromInt(payeeRealAmount)

	fee := decimal.NewFromInt(100).Sub(decimal.NewFromFloat(rate))

	amount = amount.Add(decimal.NewFromInt(singleFee)).Mul(decimal.NewFromInt(100)).Div(fee)

	// 四舍五入去除小数点
	amount = amount.Round(0)

	// 取整数
	return amount.IntPart()
}
