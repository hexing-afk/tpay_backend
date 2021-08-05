package svc

import (
	"context"
	"github.com/tal-tech/go-zero/rest"
	"log"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/merchantapi/internal/config"
	"tpay_backend/merchantapi/internal/middleware"
	"tpay_backend/pkg/cloudstorage"
	"tpay_backend/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/go-redis/redis/v8"
)

type ServiceContext struct {
	Config       config.Config
	Redis        *redis.Client
	CloudStorage cloudstorage.Storage
	DbEngine     *gorm.DB
	RedisSession *utils.RedisSession
	CheckLogin   rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	sctx := &ServiceContext{Config: c}

	// 初始化redis
	sctx.InitRedis()

	// 初始化mysql
	sctx.InitMysql()

	// 初始化文件存储
	sctx.InitCloudStorage()

	// 初始化登录session
	sctx.InitRedisSession()

	// 初始化中间件
	sctx.InitMiddleware()

	return sctx
}

// 初始化redis
func (sctx *ServiceContext) InitRedis() {
	// redis连接
	redisObj := redis.NewClient(&redis.Options{
		Addr:     sctx.Config.Redis.Host,
		Password: sctx.Config.Redis.Pass, // no password set
		DB:       sctx.Config.Redis.DB,   // use default DB
	})

	// 检测redis连接是否正常
	if err := redisObj.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("redis连接失败:%+v, err:%v", sctx.Config.Redis, err)
	}

	sctx.Redis = redisObj
}

// 初始化mysql
func (sctx *ServiceContext) InitMysql() {
	//newLogger := logger.New(
	//	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	//	logger.Config{
	//		SlowThreshold: time.Second, // Slow SQL threshold
	//		LogLevel:      logger.Info, // Log level
	//		Colorful:      false,       // Disable color
	//	},
	//)

	newLogger := utils.NewCustomGormLogger().LogMode(logger.Info)

	//启动Gorm支持
	db, err := gorm.Open(mysql.Open(sctx.Config.Mysql.DataSource), &gorm.Config{
		Logger:                 newLogger,
		SkipDefaultTransaction: true,
	})
	//如果出错就GameOver了
	if err != nil {
		log.Fatalf("mysql连接失败,err:%v,配置:%+v", err, sctx.Config.Mysql)
	}

	sctx.DbEngine = db
}

// 初始化文件存储
func (sctx *ServiceContext) InitCloudStorage() {
	var storage cloudstorage.Storage
	switch sctx.Config.CloudStorage {
	case cloudstorage.OssCloudStorage:
		storage = cloudstorage.NewOssStorage(sctx.Config.OssStorage.AccessKeyId, sctx.Config.OssStorage.SecretAccessKey,
			sctx.Config.OssStorage.Endpoint, sctx.Config.OssStorage.Bucket)
	case cloudstorage.S3CloudStorage:
		storage = cloudstorage.NewS3Storage(sctx.Config.S3Storage.AccessKeyId, sctx.Config.S3Storage.SecretAccessKey,
			sctx.Config.S3Storage.Region, sctx.Config.S3Storage.Bucket)
	}

	sctx.CloudStorage = storage
}

// 初始化登录session
func (sctx *ServiceContext) InitRedisSession() {
	config := utils.RedisSessionConfig{
		KeyPrefix: common.LoginRedisKeyPrefix,
		Expire:    common.LoginRedisExpire,
	}
	sctx.RedisSession = utils.NewRedisSession(sctx.Redis, config)
}

// 初始化中间件
func (sctx *ServiceContext) InitMiddleware() {
	sctx.CheckLogin = middleware.NewCheckLoginMiddleware().GetHandle(sctx.RedisSession)
}
