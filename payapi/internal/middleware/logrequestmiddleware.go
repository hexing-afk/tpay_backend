package middleware

import (
	"net/http"
	"time"
	"tpay_backend/utils"

	"github.com/tal-tech/go-zero/core/logx"
)

type LogRequestMiddleware struct {
}

func NewLogRequestMiddleware() *LogRequestMiddleware {
	return &LogRequestMiddleware{}
}

func (m *LogRequestMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		traceNo := utils.GetDailyId()
		begin := time.Now()

		logx.Infof("Start-TraceNo:%v|Method:%v|RequestURI:%v", traceNo, r.Method, r.RequestURI)
		//logx.Infof("Header|TraceNo:%v|Header:%v", traceNo, r.Header)
		logx.Infof("Body-TraceNo:%v|Body:%s", traceNo, string(utils.GetHttpRequestBody(r)))

		next(w, r)

		logx.Infof("End-TraceNo:%v|Duration:%v", traceNo, time.Since(begin))
	}
}
