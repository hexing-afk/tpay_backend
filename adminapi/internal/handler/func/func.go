package _func

import (
	"net/http"
	"strconv"
	"tpay_backend/adminapi/internal/common"

	"github.com/tal-tech/go-zero/core/logx"
)

func GetLoginedUserIdRequestHeader(r *http.Request) (int64, error) {
	userIdStr := r.Header.Get(common.RequestHeaderLoginedUserId)

	// 验证是否是int64
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		logx.Errorf("userId转成int64失败:err:%v", err)
		return 0, common.NewCodeError(common.UserNotLogin)
	}

	// 是否小于1
	if userId < 1 {
		logx.Errorf("读取到的userId小于1:%v", userId)
		return 0, common.NewCodeError(common.UserNotLogin)
	}

	return userId, nil
}
