#!/bin/bash
# 获取当前目录
curDir=$(cd "$(dirname "$0")"; pwd)

cd $curDir
cd ../

go build -o release/adminapid adminapi/adminapi.go
go build -o release/merchantapid merchantapi/merchantapi.go
go build -o release/payapid payapi/payapi.go

# 将"scripts"目录排除
#curDirLen=$((${#curDir}-7))

# 获取scripts目录的上级目录
#curDir=${curDir:0:${curDirLen}}

#echo $curDir

# go build -o ${curDir}release/adminapi ${curDir}adminapi/adminapi.go
