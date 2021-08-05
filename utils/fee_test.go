package utils

import "testing"

func TestCalculatePayOrderFeeMerchant(t *testing.T) {
	var reqAmount int64 = 17700000
	var singleFee int64 = 0
	rate := 4.5
	fee := CalculatePayOrderFeeMerchant(reqAmount, singleFee, rate)
	t.Logf("fee: %v", fee)
}

func TestCalculatePayOrderFeeUpstream(t *testing.T) {
	type args struct {
		reqAmount int64
		singleFee int64
		rate      float64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		// TODO: Add test cases.
		{name: "A001", args: args{reqAmount: 100, singleFee: 2, rate: 2.5}, want: 5},
		{name: "A002", args: args{reqAmount: 100, singleFee: 1, rate: 2.4}, want: 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculatePayOrderFeeUpstream(tt.args.reqAmount, tt.args.singleFee, tt.args.rate); got != tt.want {
				t.Errorf("CalculatePayOrderFeeUpstream() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateUpstreamInnerAmount(t *testing.T) {
	payeeRealAmount := int64(95)
	singleFee := int64(0)
	rate := float64(3)

	got := CalculateUpstreamInnerAmount(payeeRealAmount, singleFee, rate)
	t.Logf("got:%v", got)

}
