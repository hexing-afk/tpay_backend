package merchant

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"
	"tpay_backend/utils"

	"github.com/go-redis/redis/v8"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type EnableMerchantLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEnableMerchantLogic(ctx context.Context, svcCtx *svc.ServiceContext) EnableMerchantLogic {
	return EnableMerchantLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EnableMerchantLogic) EnableMerchant(req types.EnableMerchantRequest) error {
	var err error
	switch req.Enable {
	case model.MerchantStatusEnable:
		err = model.NewMerchantModel(l.svcCtx.DbEngine).EnableMerchant(req.MerchantId)
	case model.MerchantStatusDisable:
		err = model.NewMerchantModel(l.svcCtx.DbEngine).DisableMerchant(req.MerchantId)

		//清除商户的登录token
		l.CleanMerchantLoginToken(req.MerchantId)

	default:
		l.Errorf("操作不支持, req.Enable=%v", req.Enable)
		return common.NewCodeError(common.InvalidParam)
	}
	if err != nil {
		l.Errorf("启用|禁用商家失败, req.Enable=%v, err=%v", req.Enable, err)
		return common.NewCodeError(common.SysDBUpdate)
	}

	return nil
}

func (l *EnableMerchantLogic) CleanMerchantLoginToken(merchantId int64) {

	redisOptions := &redis.Options{
		Addr:     l.svcCtx.Config.Redis.Host,
		Password: l.svcCtx.Config.Redis.Pass, // no password set
		DB:       utils.RedisDbMerchantapi,   // merchant redis DB
	}

	// redis连接
	redisObj := redis.NewClient(redisOptions)

	// 检测redis连接是否正常
	if err := redisObj.Ping(context.Background()).Err(); err != nil {
		l.Errorf("redis连接失败:%+v, err:%v", redisOptions, err)
		return
	}

	merchantRedisSession := utils.NewRedisSession(redisObj, utils.RedisSessionConfig{
		KeyPrefix: common.MerchantLoginRedisKeyPrefix,
	})

	//清除商户的登录token
	if errDel := merchantRedisSession.CleanOtherLogined(merchantId); errDel != nil {
		l.Errorf("删除商户[%v]的登录token失败", merchantId)
		return
	}

}
