#!/bin/bash

tar -xvf /root/publish/tpay_admin_web.tar.gz --strip-components 1 -C /opt/go/tpay/admin_web/
date "+%Y-%m-%d %H:%M:%S 解压管理后台前端文件完成..."

tar -xvf /root/publish/tpay_merchant_web.tar.gz --strip-components 1 -C /opt/go/tpay/merchant_web/
date "+%Y-%m-%d %H:%M:%S 解压商户后台前端文件完成..."



rm /root/publish/tpay_admin_web.tar.gz
rm /root/publish/tpay_merchant_web.tar.gz

