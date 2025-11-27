package oss

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"io"
	"net/http"
	"path/filepath"
	"resource/internal/dao/model"
	"time"
	"toolkit/constants"
	"toolkit/errx"
	"toolkit/utils"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/zeromicro/go-zero/core/logx"
	"resource/internal/svc"
)

type UploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *UploadLogic {
	return &UploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}
func (l *UploadLogic) Upload() error {
	file, header, err := l.r.FormFile("file")
	if err != nil {
		return fmt.Errorf("failed to get form file: %w", err)
	}
	defer file.Close()

	fileName := header.Filename
	fileSuffix := filepath.Ext(fileName)
	date := time.Now().Format("2006/01/02")
	uniqueFileName := fmt.Sprintf("%s/%s%s", date, utils.GetID(), fileSuffix)

	config, err := l.getConfig()
	if err != nil {
		return fmt.Errorf("get oss config failed: %w", err)
	}

	url, err := uploadOss(l.ctx, config, uniqueFileName, file)
	if err != nil {
		return fmt.Errorf("upload to oss failed: %w", err)
	}

	ossRecord := &model.SysOss{
		OssID:        utils.GetID(),
		FileName:     uniqueFileName,
		FileSuffix:   fileSuffix,
		OriginalName: fileName,
		URL:          url,
	}
	if err := l.svcCtx.Dal.Query.SysOss.WithContext(l.ctx).Create(ossRecord); err != nil {
		return fmt.Errorf("save file record failed: %w", err)
	}
	return nil
}

func uploadOss(ctx context.Context, config *model.SysOssConfig, objectName string, reader io.Reader) (string, error) {
	provider := credentials.NewStaticCredentialsProvider(config.AccessKey, config.SecretKey)
	cfg := oss.LoadDefaultConfig().WithCredentialsProvider(provider).
		WithEndpoint(config.Endpoint).
		WithRegion(config.Region)
	client := oss.NewClient(cfg)

	exist, err := client.IsBucketExist(ctx, config.BucketName)
	if err != nil {
		return "", fmt.Errorf("failed to check bucket: %w", err)
	}
	if !exist {
		return "", fmt.Errorf("bucket %s does not exist", config.BucketName)
	}
	result, err := client.PutObject(ctx, &oss.PutObjectRequest{
		Bucket: oss.Ptr(config.BucketName),
		Key:    oss.Ptr(objectName),
		Body:   reader,
	})
	if err != nil {
		return "", fmt.Errorf("put object failed: %w", err)
	}

	fmt.Printf("Uploaded to OSS: %s (ETag: %s)\n", objectName, *result.ETag)
	url := fmt.Sprintf("%s/%s", config.Domain, objectName)
	return url, nil
}

func (l *UploadLogic) getConfig() (*model.SysOssConfig, error) {
	q := l.svcCtx.Dal.Query
	cfg := &model.SysOssConfig{}
	ex, err := l.svcCtx.Rds.ExistsCtx(l.ctx, constants.OssConfigDefaultCache)
	if !ex && err == nil {
		cfg, err = q.SysOssConfig.WithContext(l.ctx).Where(q.SysOssConfig.Status.Eq("0")).First()
		if err != nil {
			return nil, errx.GORMErrMsg(err, "未找到默认的OSS配置")
		}
		ossConfigs, err := q.SysOssConfig.WithContext(l.ctx).Find()
		if err != nil {
			return nil, err
		}
		for _, v := range ossConfigs {
			b, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			err = l.svcCtx.Rds.SetCtx(l.ctx, fmt.Sprintf(constants.OssConfigCache, v.ConfigKey), string(b))
			if err != nil {
				return nil, err
			}
			if v.Status == "0" {
				err = l.svcCtx.Rds.SetCtx(l.ctx, constants.OssConfigDefaultCache, string(b))
				if err != nil {
					return nil, err
				}
			}
		}
		return cfg, err
	}

	cfgStr, err := l.svcCtx.Rds.GetCtx(l.ctx, constants.OssConfigDefaultCache)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(cfgStr), cfg)
	return cfg, err
}
