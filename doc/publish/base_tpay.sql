/*
 Navicat Premium Data Transfer

 Source Server         : 10.41.1.242（mysql）
 Source Server Type    : MySQL
 Source Server Version : 100507
 Source Host           : 10.41.1.242:3306
 Source Schema         : test

 Target Server Type    : MySQL
 Target Server Version : 100507
 File Encoding         : 65001

 Date: 18/05/2021 16:44:14
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for admin
-- ----------------------------
DROP TABLE IF EXISTS `admin`;
CREATE TABLE `admin`  (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户名',
  `password` char(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '密码',
  `create_time` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
  `update_time` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间',
  `enable_status` tinyint UNSIGNED NOT NULL DEFAULT 1 COMMENT '账号启用状态: 1-启用， 2-禁用',
  `phone` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '手机号',
  `email` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '邮箱',
  `totp_secret` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'TOTP认证秘钥',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `username_unique`(`username`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 62 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '管理员表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of admin
-- ----------------------------
INSERT INTO `admin` VALUES (7, 'admin', 'd2f7bdc1afa7141241f868ac5d0c04dc1714e3b679f22ff52b94954daf8c6efc', 1616051322, 0, 1, '', '', 'FD6Q4TOUKEPHBI4JV2FIIF4JCXJRTISS');

-- ----------------------------
-- Table structure for admin_web_log
-- ----------------------------
DROP TABLE IF EXISTS `admin_web_log`;
CREATE TABLE `admin_web_log`  (
  `log_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '日志号',
  `admin_id` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '管理员id',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '描述',
  `type` tinyint UNSIGNED NOT NULL DEFAULT 0 COMMENT '日志类型：1-平台配置相关， 2-商家相关， ',
  `create_time` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
  PRIMARY KEY (`log_no`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of admin_web_log
-- ----------------------------

-- ----------------------------
-- Table structure for area
-- ----------------------------
DROP TABLE IF EXISTS `area`;
CREATE TABLE `area`  (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `area_name` varchar(65) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '区域名称',
  `create_time` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `area_name`(`area_name`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '地区表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of area
-- ----------------------------
INSERT INTO `area` VALUES (1, '中国', 1619005493);
INSERT INTO `area` VALUES (2, '印度', 1619005493);
INSERT INTO `area` VALUES (3, '越南', 1619005493);
INSERT INTO `area` VALUES (4, '柬埔寨', 1619005493);

-- ----------------------------
-- Table structure for currency
-- ----------------------------
DROP TABLE IF EXISTS `currency`;
CREATE TABLE `currency`  (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `currency` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '币种',
  `symbol` varchar(6) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '符号',
  `country` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '国家',
  `is_divide_hundred` tinyint UNSIGNED NOT NULL DEFAULT 1 COMMENT '前端显示是否需要除以100； 1-是, 2-否',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `currency`(`currency`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of currency
-- ----------------------------
INSERT INTO `currency` VALUES (1, 'CNY', '￥', 'China', 1);
INSERT INTO `currency` VALUES (2, 'INR', '₹', 'India', 1);
INSERT INTO `currency` VALUES (3, 'KHR', '៛', 'Cambodia', 2);
INSERT INTO `currency` VALUES (4, 'THB', '฿', 'Thailand', 1);
INSERT INTO `currency` VALUES (5, 'USD', '$', 'America', 1);
INSERT INTO `currency` VALUES (6, 'VND', '₫', 'Vietnam', 2);

-- ----------------------------
-- Table structure for global_config
-- ----------------------------
DROP TABLE IF EXISTS `global_config`;
CREATE TABLE `global_config`  (
  `config_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '全局配置key',
  `config_value` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '全局配置value',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
  `create_time` int UNSIGNED NOT NULL DEFAULT 0,
  `is_change` tinyint UNSIGNED NOT NULL DEFAULT 2 COMMENT '是否可改：1-是，2-否',
  PRIMARY KEY (`config_key`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of global_config
-- ----------------------------
INSERT INTO `global_config` VALUES ('image_base_url', 'https://tpay-dev.oss-ap-southeast-1.aliyuncs.com', '图片域名地址', 1617767519, 2);
INSERT INTO `global_config` VALUES ('payapi_host_addr', 'http://10.41.1.242:8000', 'payapi站点域名地址', 1616659140, 2);
INSERT INTO `global_config` VALUES ('pay_trade_type_slice', 'UNIONPAY,ACLEDA,UPI,ZALOPAY,MOMOPAY,VIETCOMBANK,VIETINBANKIPAY,VIETTELPAY,TPBANK,ACBBANK,usdt,UNION', '代收交易方式', 1619071669, 2);
INSERT INTO `global_config` VALUES ('site_config', '{\"SiteName\":\"\",\"SiteLogo\":\"\",\"SiteLang\":\"en-US\"}', '网站配置', 1616659140, 1);
INSERT INTO `global_config` VALUES ('totp_is_close', '1', '所有后台登录谷歌验证码是否关闭(1关闭，0开启),此开关主要是方便开发环境', 1619071679, 2);
INSERT INTO `global_config` VALUES ('transfer_trade_type_slice', 'TRANSFER_UNIONPAY,TRANSFER_ACLEDA,TRANSFER_ZALO,TRANSFER_MOMO,TRANSFER_BANK,TRANSFER_UNION', '代付交易方式', 1619071679, 2);

-- ----------------------------
-- Table structure for merchant
-- ----------------------------
DROP TABLE IF EXISTS `merchant`;
CREATE TABLE `merchant`  (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `merchant_no` char(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '商户编号',
  `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户名',
  `password` char(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '密码',
  `create_time` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
  `update_time` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间',
  `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '手机号',
  `email` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '邮箱',
  `currency` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '币种(来源于currency表)',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '账号状态：1-启用, 2-禁用',
  `md5_key` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'md5通信秘钥',
  `balance` bigint NOT NULL DEFAULT 0 COMMENT '余额',
  `frozen_amount` bigint UNSIGNED NOT NULL DEFAULT 0 COMMENT '冻结金额',
  `pay_password` char(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '支付密码',
  `ip_white_list` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT 'ip白名单(英文逗号分隔)',
  `totp_secret` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'TOTP认证秘钥',
  `area_id` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '地区id(来自area表id)',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `username_unique`(`username`) USING BTREE COMMENT '商家账号唯一',
  UNIQUE INDEX `merchant_no_unique`(`merchant_no`) USING BTREE COMMENT '商家编号唯一'
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '商户表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of merchant
-- ----------------------------

-- ----------------------------
-- Table structure for merchant_bank_card
-- ----------------------------
DROP TABLE IF EXISTS `merchant_bank_card`;
CREATE TABLE `merchant_bank_card`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `merchant_id` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '商户id',
  `create_time` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
  `bank_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '银行名称',
  `bank_code` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '银行代码',
  `branch_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '支行名称',
  `account_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '开户名',
  `card_number` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '银行卡号',
  `currency` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '币种',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
  `update_time` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `bank_name_card_number_unique`(`merchant_id`, `bank_name`, `card_number`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of merchant_bank_card
-- ----------------------------

-- ----------------------------
-- Table structure for merchant_channel
-- ----------------------------
DROP TABLE IF EXISTS `merchant_channel`;
CREATE TABLE `merchant_channel`  (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `merchant_id` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '商户Id',
  `platform_channel_id` int NOT NULL DEFAULT 0 COMMENT '平台通道id',
  `rate` float(10, 2) UNSIGNED NOT NULL DEFAULT 0 COMMENT '商户通道费率; 保留两位小数点',
  `single_fee` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '单笔手续费',
  `status` tinyint UNSIGNED NOT NULL DEFAULT 0 COMMENT '1-启用, 2-禁用',
  `create_time` int UNSIGNED NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `merchant_channel_key`(`merchant_id`, `platform_channel_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of merchant_channel
-- ----------------------------

-- ----------------------------
-- Table structure for merchant_channel_upstream
-- ----------------------------
DROP TABLE IF EXISTS `merchant_channel_upstream`;
CREATE TABLE `merchant_channel_upstream`  (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `merchant_id` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '商户id',
  `merchant_channel_id` int NOT NULL DEFAULT 0 COMMENT '商户通道id',
  `upstream_channel_id` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '上游通道id',
  `weight` mediumint UNSIGNED NOT NULL DEFAULT 10 COMMENT '权重',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of merchant_channel_upstream
-- ----------------------------

-- ----------------------------
-- Table structure for merchant_recharge_order
-- ----------------------------
DROP TABLE IF EXISTS `merchant_recharge_order`;
CREATE TABLE `merchant_recharge_order`  (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '商户充值记录id',
  `order_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '订单号',
  `order_amount` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '订单金额',
  `merchant_id` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '商户id',
  `order_status` tinyint UNSIGNED NOT NULL DEFAULT 0 COMMENT '订单状态:1-待处理；2-通过； 3-驳回',
  `create_time` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
  `recharge_remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '充值备注',
  `platform_bank_card_id` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '平台收款卡id',
  `finish_time` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '处理完成时间',
  `audit_remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '审核备注',
  `bank_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '收款银行',
  `card_number` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '收款卡号',
  `payee_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '收款人姓名',
  `branch_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '支行名称',
  `currency` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '币种',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `order_no`(`order_no`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of merchant_recharge_order
-- ----------------------------

-- ----------------------------
-- Table structure for merchant_wallet_log
-- ----------------------------
DROP TABLE IF EXISTS `merchant_wallet_log`;
CREATE TABLE `merchant_wallet_log`  (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `merchant_id` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '商户id',
  `create_time` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
  `op_type` tinyint NOT NULL DEFAULT 0 COMMENT '变动类型：1+，2-',
  `change_amount` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '变动金额',
  `after_balance` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '变动后余额',
  `business_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '业务单号:(充值；提现；代付；收款)',
  `source` tinyint UNSIGNED NOT NULL DEFAULT 0 COMMENT '变动来源：1-手动调账；2-充值；3-提现；4-代付；5-收款；',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of merchant_wallet_log
-- ----------------------------

-- ----------------------------
-- Table structure for merchant_withdraw_config
-- ----------------------------
DROP TABLE IF EXISTS `merchant_withdraw_config`;
CREATE TABLE `merchant_withdraw_config`  (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `merchant_id` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '商户id',
  `deduction_method` tinyint UNSIGNED NOT NULL DEFAULT 0 COMMENT '手续费扣费方式: 1-内扣， 2-外扣',
  `single_min_amount` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '单笔提现最小金额',
  `single_max_amount` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '单笔提现最大金额',
  `rate` float(10, 2) UNSIGNED NOT NULL DEFAULT 0 COMMENT '提现费率; 保留两位小数点',
  `single_fee` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '单笔提现手续费',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `merchant_id_key`(`merchant_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of merchant_withdraw_config
-- ----------------------------

-- ----------------------------
-- Table structure for merchant_withdraw_order
-- ----------------------------
DROP TABLE IF EXISTS `merchant_withdraw_order`;
CREATE TABLE `merchant_withdraw_order`  (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '商户提现订单id',
  `order_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '订单号',
  `merchant_id` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '商户Id',
  `order_amount` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '订单金额',
  `merchant_fee` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '商户手续费',
  `real_amount` int NOT NULL DEFAULT 0 COMMENT '实际到账金额',
  `decrease_amount` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '扣减商户账号金额',
  `currency` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '币种',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '提现备注',
  `bank_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '收款银行',
  `bank_code` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '银行代码',
  `branch_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '支行名称',
  `payee_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '收款人',
  `card_number` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '收款卡号',
  `create_time` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
  `order_status` tinyint UNSIGNED NOT NULL DEFAULT 0 COMMENT '订单状态: 1-待处理; 2-驳回; 3-通过审核； 4-派单中； 5-派单成功； 6-派单失败；7-成功',
  `audit_time` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '审核时间',
  `audit_remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '审核备注',
  `merchant_bank_card_id` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '商户银行卡id',
  `deduction_method` tinyint UNSIGNED NOT NULL DEFAULT 0 COMMENT '商户手续费扣费方式: 1-内扣; 2-外扣; 来源于商户提现配置表',
  `success_remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '付款成功备注',
  `transfer_order_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '代付单号(派单单号)',
  `area_id` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '地区id(来自area表id)',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `order_no`(`order_no`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of merchant_withdraw_order
-- ----------------------------

-- ----------------------------
-- Table structure for order_notify_log
-- ----------------------------
DROP TABLE IF EXISTS `order_notify_log`;
CREATE TABLE `order_notify_log`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `order_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '内部订单号',
  `merchant_order_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '商户订单号',
  `status` tinyint UNSIGNED NOT NULL DEFAULT 0 COMMENT '通知状态(0通知中,1成功,2失败)',
  `result` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '通知返回的结果',
  `create_time` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '通知时间',
  `order_type` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '订单类型(1代收订单, 2代付订单)',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of order_notify_log
-- ----------------------------

-- ----------------------------
-- Table structure for pay_order
-- ----------------------------
DROP TABLE IF EXISTS `pay_order`;
CREATE TABLE `pay_order`  (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `order_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '平台单号',
  `merchant_order_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '商户订单号(即下游订单号)',
  `upstream_order_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '上游订单号',
  `merchant_no` char(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '商户编号',
  `req_amount` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '订单请求金额',
  `increase_amount` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '账户增加的金额',
  `merchant_fee` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '商户手续费',
  `upstream_amount` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '请求上游的金额',
  `upstream_fee` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '上游手续费',
  `order_status` tinyint UNSIGNED NOT NULL DEFAULT 0 COMMENT '订单状态:1-待支付;2-已支付;3-支付失败',
  `currency` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '币种',
  `create_time` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
  `update_time` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间',
  `notify_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '异步通知url',
  `return_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '同步跳转url',
  `platform_channel_id` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '平台通道id',
  `upstream_channel_id` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '上游通道id',
  `notify_status` tinyint UNSIGNED NOT NULL DEFAULT 0 COMMENT '异步通知状态(0未通知,1成功,2通知进行中,3超时)',
  `notify_fail_times` tinyint UNSIGNED NOT NULL DEFAULT 0 COMMENT '通知失败的次数',
  `next_notify_time` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '下次通知时间(通知失败时才有)',
  `subject` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '商品的标题/交易标题/订单标题/订单关键字等',
  `area_id` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '地区id(来自area表id)',
  `mode` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '模式：test|pro(测试|生产)',
  `merchant_rate` float(10, 2) UNSIGNED NOT NULL DEFAULT 0 COMMENT '商户费率',
  `merchant_single_fee` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '商户单笔手续费',
  `payment_amount` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '实际支付金额',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `merchant_order_unique`(`merchant_no`, `merchant_order_no`) USING BTREE COMMENT '同一商户的订单只能有一个',
  UNIQUE INDEX `order_no`(`order_no`) USING BTREE COMMENT '平台订单号唯一'
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '代收订单表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of pay_order
-- ----------------------------

-- ----------------------------
-- Table structure for platform_bank_card
-- ----------------------------
DROP TABLE IF EXISTS `platform_bank_card`;
CREATE TABLE `platform_bank_card`  (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'id',
  `bank_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '银行名称',
  `account_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '开户名',
  `create_time` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
  `card_number` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '银行卡号',
  `branch_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '支行名称',
  `currency` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '币种',
  `max_amount` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '最大收款额度',
  `qr_code` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '收款二维码',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '状态：1-启用， 2-禁用',
  `today_received` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '今日已收金额',
  `bank_code` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '银行代码',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `bank_name_card_number_currency_unique`(`bank_name`, `card_number`, `currency`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of platform_bank_card
-- ----------------------------

-- ----------------------------
-- Table structure for platform_channel
-- ----------------------------
DROP TABLE IF EXISTS `platform_channel`;
CREATE TABLE `platform_channel`  (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '平台通道id',
  `channel_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '通道名称',
  `channel_code` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '通道代码',
  `create_time` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
  `channel_desc` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '通道描述',
  `channel_type` tinyint NOT NULL DEFAULT 0 COMMENT '通道类型:1-代收, 2-代付',
  `update_time` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间',
  `status` tinyint UNSIGNED NOT NULL DEFAULT 0 COMMENT '状态: 1-启用， 2-禁用',
  `area_id` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '地区id(来自area表id)',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `channel_name_unique`(`channel_name`) USING BTREE,
  UNIQUE INDEX `channel_code_type_unique`(`channel_code`, `channel_type`, `area_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of platform_channel
-- ----------------------------

-- ----------------------------
-- Table structure for platform_channel_upstream
-- ----------------------------
DROP TABLE IF EXISTS `platform_channel_upstream`;
CREATE TABLE `platform_channel_upstream`  (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `upstream_channel_id` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '上游通道id',
  `platform_channel_id` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '平台通道id',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `upstrean_downstream_channel_key`(`upstream_channel_id`, `platform_channel_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of platform_channel_upstream
-- ----------------------------

-- ----------------------------
-- Table structure for platform_wallet_log
-- ----------------------------
DROP TABLE IF EXISTS `platform_wallet_log`;
CREATE TABLE `platform_wallet_log`  (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '平台钱包日志id',
  `business_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '业务号',
  `source` tinyint UNSIGNED NOT NULL DEFAULT 0 COMMENT '收益来源:3-商户提现；4-商户代付；5-商户代收；',
  `merchant_fee` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '商户手续费',
  `upstream_fee` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '上游手续费',
  `income` int NOT NULL DEFAULT 0 COMMENT '收益',
  `create_time` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
  `currency` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '币种',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of platform_wallet_log
-- ----------------------------

-- ----------------------------
-- Table structure for transfer_order
-- ----------------------------
DROP TABLE IF EXISTS `transfer_order`;
CREATE TABLE `transfer_order` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `order_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '平台订单号',
  `merchant_order_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '商户订单号',
  `upstream_order_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '上游订单号',
  `merchant_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '商户编号',
  `create_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `req_amount` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '订单请求金额',
  `merchant_fee` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商户手续费',
  `decrease_amount` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '账户扣除的金额',
  `upstream_amount` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '请求上游的金额',
  `upstream_fee` int(10) NOT NULL COMMENT '上游手续费',
  `currency` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '币种',
  `order_status` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '订单状态: 1-待支付; 2-支付成功; 3-支付失败',
  `platform_channel_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '平台通道id',
  `upstream_channel_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '上游通道id',
  `update_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `notify_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '异步通知地址',
  `return_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '同步跳转地址',
  `bank_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '收款银行名称',
  `account_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '银行卡开户名',
  `card_number` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '收款卡号',
  `branch_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '支行名称',
  `notify_status` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '异步通知状态(0未通知,1成功,2通知进行中,3超时)',
  `notify_fail_times` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '通知失败次数',
  `next_notify_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '下次通知时间',
  `payee_real_amount` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '收款方实际到账金额',
  `fee_deduct_type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '手续费扣款方式(1内扣,2外扣)',
  `upstream_fail_reason` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '上游失败原因',
  `order_source` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '代付订单来源;1-接口; 2-平台提现派单；3-商户后台付款',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '付款备注',
  `bank_code` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '银行代码(ifsc_code)',
  `area_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '地区id(来自area表id)',
  `mode` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '模式：test|pro(测试|生产)',
  `merchant_rate` float(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '商户费率',
  `merchant_single_fee` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商户单笔手续费',
  `batch_no` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `batch_row_no` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `merchant_order_no_unique` (`merchant_order_no`,`merchant_no`) USING BTREE,
  UNIQUE KEY `order_no_unique` (`order_no`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='代付订单表';
-- ----------------------------
-- Records of transfer_order
-- ----------------------------

-- ----------------------------
-- Table structure for upstream
-- ----------------------------
DROP TABLE IF EXISTS `upstream`;
CREATE TABLE `upstream`  (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `upstream_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '上游名称',
  `upstream_merchant_no` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '上游账号id',
  `upstream_code` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '上游代码',
  `call_config` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL COMMENT '通信配置',
  `create_time` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
  `area_id` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '地区id(来自area表id)',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `upstream_merchant_no_idx`(`upstream_merchant_no`) USING BTREE,
  UNIQUE INDEX `upstream_name`(`upstream_name`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '上游表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of upstream
-- ----------------------------

-- ----------------------------
-- Table structure for upstream_channel
-- ----------------------------
DROP TABLE IF EXISTS `upstream_channel`;
CREATE TABLE `upstream_channel`  (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `channel_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '通道名称',
  `create_time` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
  `channel_code` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '通道代码',
  `channel_desc` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '通道描述',
  `status` tinyint UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态: 1-启用， 2-禁用',
  `currency` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '通道币种',
  `upstream_id` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '上游id',
  `channel_type` tinyint UNSIGNED NOT NULL DEFAULT 0 COMMENT '通道类型:1-代收, 2-代付',
  `update_time` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间',
  `deduction_method` tinyint NOT NULL DEFAULT 0 COMMENT '扣费方式: 1-内扣， 2-外扣',
  `rate` float(10, 2) UNSIGNED NOT NULL DEFAULT 0 COMMENT '通道费率; 保留两位小数点(单位:百分比)',
  `single_fee` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '单笔手续费',
  `single_max_amount` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '单笔最大金额',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `channel_name_unique`(`channel_name`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '上游通道表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of upstream_channel
-- ----------------------------

SET FOREIGN_KEY_CHECKS = 1;
