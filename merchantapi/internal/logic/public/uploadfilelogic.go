package public

import (
	"context"
	"fmt"
	"mime/multipart"
	"path"
	"time"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/model"
	"tpay_backend/pkg/cloudstorage"

	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type UploadFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) UploadFileLogic {
	return UploadFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadFileLogic) UploadFile(merchantId int64, fileHeader *multipart.FileHeader) (*types.UploadFileReply, error) {

	filenameWithSuffix := path.Ext(fileHeader.Filename)
	l.Infof("上传的文件名[%v], 后缀名[%v]", fileHeader.Filename, filenameWithSuffix)

	baseName := fmt.Sprintf("%d%s", time.Now().Unix(), filenameWithSuffix)
	fileName := ""

	//转存到该存的地方
	switch filenameWithSuffix {
	case ".xlsx":
		fileName = path.Join(cloudstorage.Xlsx_Dir, baseName)
		//fileName = path.Join(cloudstorage.Misc_Dir, fileHeader.Filename)
	default:
		l.Errorf("文件后缀名[%v]错误", filenameWithSuffix)
		return nil, common.NewCodeError(common.NotSupportUploadFile)
	}

	//
	l.Infof("开始上传文件[%v]到云存储", fileName)

	if err := l.svcCtx.CloudStorage.UploadByMultipartFileHeader(fileName, fileHeader, true); err != nil {
		l.Errorf("上传到云存储失败，err:[%v]", err)
		return nil, common.NewCodeError(common.UploadFail)
	}

	l.Infof("上传文件[%v]到云存储结束", fileName)

	//添加上传日志
	if err := model.NewUploadFileLogModel(l.svcCtx.DbEngine).Insert(&model.UploadFileLog{
		FileName:    fileName,
		AccountId:   merchantId,
		AccountType: model.UploadFileLogAccountTypeMerchant,
	}); err != nil {
		l.Errorf("插入日志失败，err:[%v]", err)
		return nil, common.NewCodeError(common.UploadFail)
	}

	return &types.UploadFileReply{
		FileName: fileName,
	}, nil
}
