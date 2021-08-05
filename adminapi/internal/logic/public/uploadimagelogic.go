package logic

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"
	"tpay_backend/utils"

	"github.com/tal-tech/go-zero/core/logx"
)

type UploadImageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) UploadImageLogic {
	return UploadImageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadImageLogic) UploadImage(req types.UploadImageRequest) (*types.UploadImageResponse, error) {
	//处理图片
	imageData, err := utils.DealWithImage(req.ImageStr, "")
	if err != nil {
		l.Logger.Errorf("组装图片失败, err=%v", err)
		return nil, common.NewCodeError(common.UploadFail)
	}

	logx.Info("开始上传文件到云存储")

	err = l.svcCtx.CloudStorage.UploadByContent(imageData.FileName, imageData.ImageContent, true)
	if err != nil {
		l.Logger.Errorf("上传文件失败, err:%v, FileName: %s", err, imageData.FileName)
		return nil, common.NewCodeError(common.UploadFail)
	}
	logx.Infof("上传文件到云存储结束, FileName: %s", imageData.FileName)

	return &types.UploadImageResponse{
		ImageUrl: imageData.FileName,
	}, nil
}
