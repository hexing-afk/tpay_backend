package middleware

import (
	"fmt"
	"net/http"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/utils"

	"github.com/tal-tech/go-zero/core/logx"

	"github.com/tal-tech/go-zero/rest/httpx"
)

type CheckLoginMiddleware struct {
}

func NewCheckLoginMiddleware() *CheckLoginMiddleware {
	return &CheckLoginMiddleware{}
}

func (m *CheckLoginMiddleware) GetHandle(redisSession *utils.RedisSession) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			loginToken := r.Header.Get(common.LoginTokenHeaderKey)
			if loginToken == "" {
				httpx.Error(w, common.NewCodeError(common.NoLoginToken))
				return
			}

			// 解析登录token
			userId, token, err := common.LoginTokenParse(loginToken)
			if err != nil {
				logx.Errorf("解析用户登录token失败:err[%v]", err)
				httpx.Error(w, common.NewCodeError(common.LoginTokenParseFailed))
				return
			}

			// 检查是否已在redis登录
			if err := redisSession.IsLogined(userId, token); err != nil {
				logx.Errorf("检查用户是否已在redis登录失败:err:%v", err)
				httpx.Error(w, common.NewCodeError(common.UserNotLogin))
				return
			}

			// 刷新用户登录token有效期
			if err := redisSession.Refresh(userId, token); err != nil {
				logx.Errorf("刷新用户登录token有效期失败:%v", err)
			}

			// 设置到request的头部，以便在handler中使用
			r.Header.Set(common.RequestHeaderLoginedUserId, fmt.Sprintf("%v", userId))

			// 设置到request的头部，以便在handler中使用
			r.Header.Set(common.RequestHeaderLoginedToken, token)

			next(w, r)
		}
	}
}
