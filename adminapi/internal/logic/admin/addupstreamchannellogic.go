package admin

import (
	"context"
	"strings"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type AddUpstreamChannelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddUpstreamChannelLogic(ctx context.Context, svcCtx *svc.ServiceContext) AddUpstreamChannelLogic {
	return AddUpstreamChannelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddUpstreamChannelLogic) AddUpstreamChannel(req types.AddUpstreamChannelRequest) error {
	data := &model.UpstreamChannel{
		ChannelName: req.ChannelName,
		ChannelCode: req.ChannelCode,
		ChannelDesc: req.ChannelDesc,
		Status:      model.UpstreamChannelStatusEnable,
		Currency:    req.Currency,
		UpstreamId:  req.UpstreamId,
		ChannelType: req.ChannelType,
	}
	if err := model.NewUpstreamChannelModel(l.svcCtx.DbEngine).Insert(data); err != nil {
		l.Errorf("添加上游通道[%v]失败, err=%v", req.ChannelName, err)
		errStr := err.Error()
		if strings.Contains(errStr, model.UpstreamChannelNameUniqueErr.Error()) || strings.Contains(errStr, model.UpstreamChannelCodeUniqueErr.Error()) {
			return common.NewCodeError(common.ChannelRepetition)
		} else {
			return common.NewCodeError(common.SysDBAdd)
		}
	}

	return nil
}
