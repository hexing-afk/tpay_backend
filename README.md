
#### 项目列表
1. adminapi 平台管理后台
2. merchantapi 商户管理后台
3. payapi api入口


#### golang版本: go1.16.3 
#### golang框架: go-zero v1.1.5

#### 命令行自动生成模型文件 
```shell
$ goctl model mysql datasource -url="root:123456@tcp(127.0.0.1:3306)/tpay" -table="admin" -dir .
```

#### adminapi生成接口相关文件 
```shell
$ goctl api go -api adminapi.api -dir .
```


#### merchantapi生成接口相关文件 
```shell
$ goctl api go -api merchantapi.api -dir .
```

#### payapi生成接口相关文件 
```shell
$ goctl api go -api payapi.api -dir .
```