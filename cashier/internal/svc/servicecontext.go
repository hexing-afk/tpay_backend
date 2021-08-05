package svc

import (
	"context"
	"log"
	"tpay_backend/cashier/internal/config"
	"tpay_backend/utils"

	"gorm.io/gorm/logger"

	"gorm.io/driver/mysql"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config   config.Config
	Redis    *redis.Client
	DbEngine *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	sctx := &ServiceContext{Config: c}

	// 初始化mysql
	sctx.InitMysql()
	// 初始化redis
	sctx.InitRedis()

	return sctx
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
