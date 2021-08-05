
-- 只保留admin账号
DELETE FROM `admin` WHERE username != 'admin';
TRUNCATE TABLE `admin_web_log`;
TRUNCATE TABLE `merchant`;
TRUNCATE TABLE `merchant_bank_card`;
TRUNCATE TABLE `merchant_channel`;
TRUNCATE TABLE `merchant_channel_upstream`;
TRUNCATE TABLE `merchant_recharge_order`;
TRUNCATE TABLE `merchant_wallet_log`;
TRUNCATE TABLE `merchant_withdraw_config`;
TRUNCATE TABLE `merchant_withdraw_order`;
TRUNCATE TABLE `order_notify_log`;
TRUNCATE TABLE `pay_order`;
TRUNCATE TABLE `platform_bank_card`;
TRUNCATE TABLE `platform_channel`;
TRUNCATE TABLE `platform_channel_upstream`;
TRUNCATE TABLE `platform_wallet_log`;
TRUNCATE TABLE `transfer_order`;
TRUNCATE TABLE `upstream`;
TRUNCATE TABLE `upstream_channel`;