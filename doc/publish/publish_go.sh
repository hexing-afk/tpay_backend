#!/bin/bash

serviceList="adminapi merchantapi payapi"

chmod 755  /root/publish/*.upload

for service in ${serviceList}; do
	mv /root/publish/${service}d.upload /opt/go/xpay/${service}d
	date "+%Y-%m-%d %H:%M:%S 复制${service}文件"
done
	
date "+%Y-%m-%d %H:%M:%S 全部复制完成..."


for service in ${serviceList}; do
	supervisorctl restart ${service}
	date "+%Y-%m-%d %H:%M:%S 重启${service}服务完成"
done


date "+%Y-%m-%d %H:%M:%S 执行完成..."

supervisorctl status
