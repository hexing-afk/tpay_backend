package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/tal-tech/go-zero/core/logx"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type CustomGormLogger struct {
	SlowThreshold time.Duration
	Colorful      bool
	LogLevel      logger.LogLevel
	infoStr       string
	warnStr       string
	errStr        string
	traceStr      string
	traceErrStr   string
	traceWarnStr  string
}

func NewCustomGormLogger() logger.Interface {
	l := &CustomGormLogger{}
	l.SlowThreshold = time.Second * 3

	l.infoStr = "[gorm-info] file:%s "
	l.warnStr = "[gorm-warn] file:%s "
	l.errStr = "[gorm-error] file:%s "

	l.traceStr = "[gorm-trace] file:%s [%.3fms] [rows:%v] %s"
	//l.traceStr = "[gorm-trace] file:%s duration:[%.3fms] [rows:%v] sql:%s"

	// 大于SlowThreshold的值会有traceWarnStr
	l.traceWarnStr = "[gorm-trace-warn] file:%s %s [%.3fms] [rows:%v] %s"
	//l.traceWarnStr = "[gorm-trace-warn] file:%s slowLog:%s duration:[%.3fms] [rows:%v] sql:%s"

	// sql语句错误会有traceErrStr
	l.traceErrStr = "[gorm-trace-error] file:%s %s [%.3fms] [rows:%v] %s"
	//l.traceWarnStr = "[gorm-trace-warn] file:%s error:%s duration:[%.3fms] [rows:%v] sql:%s"
	return l
}

func (l *CustomGormLogger) Printf(s string, data ...interface{}) {
	logx.Infof(s, data...)
}

// LogMode log mode
func (l *CustomGormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l CustomGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		logx.Infof(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l CustomGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		logx.Errorf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l CustomGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		logx.Errorf(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
func (l CustomGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel > logger.Silent {
		elapsed := time.Since(begin)
		switch {
		case err != nil && l.LogLevel >= logger.Error:
			sql, rows := fc()
			if rows == -1 {
				logx.Errorf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				logx.Errorf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
			if rows == -1 {
				logx.Errorf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				logx.Errorf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case l.LogLevel == logger.Info:
			sql, rows := fc()
			if rows == -1 {
				logx.Infof(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				logx.Infof(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		}
	}
}
