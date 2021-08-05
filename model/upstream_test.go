package model

import (
	"testing"
	"tpay_backend/test"
)

func TestUpstreamModel_FindOneByUpstreamMerchantNo(t *testing.T) {
	upstreamModel := NewUpstreamModel(test.DbEngine)

	upMerchantNo := "3d2f9bb9-f370-412e-91fb-d011f76706f3"

	got, err := upstreamModel.FindOneByUpstreamMerchantNo(upMerchantNo)

	t.Logf("err:%v", err)
	t.Logf("got:%v", got)
}
