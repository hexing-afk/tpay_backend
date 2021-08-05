package _func

import (
	"github.com/tal-tech/go-zero/core/logx"
	"net/http"
	"strconv"
	"tpay_backend/merchantapi/internal/common"
)

func GetLoginedUserIdRequestHeader(r *http.Request) (int64, error) {
	userIdStr := r.Header.Get(common.RequestHeaderLoginedUserId)

	// 验证是否是int64
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		logx.Errorf("userId转int64失败err:%v", err)
		return 0, common.NewCodeError(common.UserNotLogin)
	}

	// 是否小于1
	if userId < 1 {
		logx.Errorf("userId小于1:%v", userId)
		return 0, common.NewCodeError(common.UserNotLogin)
	}

	return userId, nil
}
