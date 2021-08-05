package test

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DbEngine *gorm.DB

func init() {
	dataSource := "root:123456@tcp(127.0.0.1:3306)/tpay?charset=utf8mb4"

	//启动Gorm支持
	db, err := gorm.Open(mysql.Open(dataSource), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	//如果出错就GameOver了
	if err != nil {
		log.Fatalf("mysql连接失败,err:%v,配置:%+v", err, dataSource)
	}

	DbEngine = db
}
