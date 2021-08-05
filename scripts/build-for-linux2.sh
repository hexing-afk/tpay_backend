#!/bin/bash
# 获取当前目录
curDir=$(cd "$(dirname "$0")"; pwd)

cd $curDir
cd ../

# 服务列表
appList="adminapi merchantapi payapi cashier"

for app in $appList; do
  go build -o release/${app}d ${app}/${app}.go
  echo build ${app} success
done