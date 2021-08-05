package model

import (
	"testing"
	"tpay_backend/test"
)

func TestPlatformChannelUpstreamModel_FindUpChannelById(t *testing.T) {

	m := NewPlatformChannelModel(test.DbEngine)

	var id int64 = 99

	data, err := m.FindUpChannelById(id)
	if err != nil {
		t.Errorf("FindUpChannelById() error = %v", err)
		return
	}
	t.Logf("data: %+v", data)
}
