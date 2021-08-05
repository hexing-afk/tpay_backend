package public

import (
	"context"
	"errors"
	"fmt"
	"github.com/tal-tech/go-zero/core/logx"
	"strings"
	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/model"
)

// 代付API路径
const PayApiTransferPath = "/system/transfer"

// 批量代付API路径
const PayApiBatchTransferPath = "/system/transfer-batch"

type FuncLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFuncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FuncLogic {
	return &FuncLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 获取代付接口地址
func (l *FuncLogic) SystemTransferApiUrl() (string, error) {
	host, err := model.NewGlobalConfigModel(l.svcCtx.DbEngine).FindValueByKey(model.ConfigPayapiHostAddr)
	if err != nil {
		return "", errors.New(fmt.Sprintf("查询全局配置[%v]失败：%v", model.ConfigPayapiHostAddr, err))
	}

	if strings.TrimSpace(host) == "" {
		l.Errorf("", model.ConfigPayapiHostAddr, err)
		return "", errors.New(fmt.Sprintf("系统没有配全局配置[%v]", model.ConfigPayapiHostAddr))
	}

	return strings.TrimRight(host, "/") + PayApiTransferPath, nil
}

// 获取批量代付接口地址
func (l *FuncLogic) SystemTransferBatchApiUrl() (string, error) {
	host, err := model.NewGlobalConfigModel(l.svcCtx.DbEngine).FindValueByKey(model.ConfigPayapiHostAddr)
	if err != nil {
		return "", errors.New(fmt.Sprintf("查询全局配置[%v]失败：%v", model.ConfigPayapiHostAddr, err))
	}

	if strings.TrimSpace(host) == "" {
		l.Errorf("", model.ConfigPayapiHostAddr, err)
		return "", errors.New(fmt.Sprintf("系统没有配全局配置[%v]", model.ConfigPayapiHostAddr))
	}

	return strings.TrimRight(host, "/") + PayApiBatchTransferPath, nil
}
