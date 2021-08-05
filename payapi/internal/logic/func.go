package logic

import (
	"errors"
	"fmt"
	"strings"
	"tpay_backend/model"
	"tpay_backend/payapi/internal/svc"
	"tpay_backend/upstream"

	"github.com/tal-tech/go-zero/core/logx"
)

type FuncLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewFuncLogic(svcCtx *svc.ServiceContext) *FuncLogic {
	return &FuncLogic{
		svcCtx: svcCtx,
	}
}

// 获取上游
func (l *FuncLogic) GetUpstream(upstreamChannelId int64) (upstream.Upstream, error) {
	var err error
	// 1.获取对应上游
	up, err := model.NewUpstreamChannelModel(l.svcCtx.DbEngine).FindUpstreamByChannelId(upstreamChannelId)
	if err != nil {
		if err == model.ErrRecordNotFound {
			return nil, errors.New(fmt.Sprintf("未找到对应的上游:upstreamChannelId:%v", upstreamChannelId))
		} else {
			return nil, errors.New(fmt.Sprintf("查询上游信息失败:err:%v,upstreamChannelId:%v", err, upstreamChannelId))
		}
	}

	obj, err := l.GetUpstreamObject(up)

	return obj, err
}

// 获取对应对象
func (l *FuncLogic) GetUpstreamObject(up *model.Upstream) (upstream.Upstream, error) {
	var upstreamObj upstream.Upstream
	var err error

	// 2.初始化上游
	switch up.UpstreamCode {
	case upstream.UpstreamCodeTotopay: // totopay
		upstreamObj, err = upstream.NewTotopay(up.UpstreamMerchantNo, up.CallConfig)
	case upstream.UpstreamCodeGoldPays:
		upstreamObj, err = upstream.NewGoldPays(up.UpstreamMerchantNo, up.CallConfig)
	case upstream.UpstreamCodeZf777Pay:
		upstreamObj, err = upstream.NewThreeSevenPay(up.UpstreamMerchantNo, up.CallConfig)
	case upstream.UpstreamCodeXPay:
		upstreamObj, err = upstream.NewXPay(up.UpstreamMerchantNo, up.CallConfig)
	case upstream.UpstreamCodeToppay:
		upstreamObj, err = upstream.NewTopPay(up.UpstreamMerchantNo, up.CallConfig)
	}

	if err != nil {
		return nil, errors.New(fmt.Sprintf("初始化上游失败err:%v,upstream:%v,config:%v", err, up.UpstreamName, up.CallConfig))
	}

	if upstreamObj == nil {
		return nil, errors.New(fmt.Sprintf("上游(%v)未配置", up.UpstreamName))
	}

	return upstreamObj, nil
}

// 获取上游异步通知地址
func (l *FuncLogic) GetUpstreamNotifyUrl(notifyPath string) (string, error) {
	// 获取payapi站点域名
	host, err := model.NewGlobalConfigModel(l.svcCtx.DbEngine).FindValueByKey(model.ConfigPayapiHostAddr)
	if err != nil {
		return "", errors.New(fmt.Sprintf("查询网站配置失败,key:%v,err=%v", model.ConfigPayapiHostAddr, err))
	}

	if strings.TrimSpace(host) == "" {
		return "", errors.New("系统没有配置异步回调地址")
	}

	if strings.TrimSpace(notifyPath) == "" {
		return "", errors.New("系统没有配置异步回调地址路径")
	}

	notifyPath = strings.TrimPrefix(notifyPath, "/")
	notifyUrl := strings.TrimRight(host, "/") + "/" + strings.TrimRight(notifyPath, "/")

	return notifyUrl, nil
}
