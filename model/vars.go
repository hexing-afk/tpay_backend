package model

import (
	"errors"
	"gorm.io/gorm"
)

var (
	ErrRecordNotFound            = gorm.ErrRecordNotFound
	AdminUserUniqueErr           = errors.New("for key 'username_unique'")
	PlatformChannelNameUniqueErr = errors.New("for key 'channel_name_unique'")
	PlatformChannelCodeUniqueErr = errors.New("for key 'channel_code_unique'")
	UpstreamChannelNameUniqueErr = errors.New("for key 'channel_name_unique'")
	UpstreamChannelCodeUniqueErr = errors.New("for key 'channel_code_unique'")
	BalanceErr                   = errors.New("Error 1690: BIGINT UNSIGNED value is out of range in '`tpay`.`merchant`.`balance` - ")
)
