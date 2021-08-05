package public

import (
	"context"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/model"

	"github.com/tal-tech/go-zero/core/logx"
	"tpay_backend/merchantapi/internal/svc"
)

type BatchTransferFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchTransferFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) BatchTransferFileLogic {
	return BatchTransferFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchTransferFileLogic) BatchTransferFile() ([]byte, error) {
	fileName, err := model.NewGlobalConfigModel(l.svcCtx.DbEngine).BaseBatchTransferFileName()
	if err != nil {
		l.Errorf("查询不到批量付款文件名称, err:%v", err)
		return nil, err
	}

	//从云存储获取文件
	bytes, err := l.svcCtx.CloudStorage.GetObject(fileName)
	if err != nil {
		l.Errorf("从云存储获取文件失败,FileName:%s, err:%v", fileName, err)
		return nil, common.NewCodeError(common.FileNotExist)
	}

	return bytes, nil
}
