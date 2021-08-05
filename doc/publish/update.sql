SET FOREIGN_KEY_CHECKS = 0;

    CREATE TABLE `transfer_batch_order`  (
      `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
      `batch_no` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '批量号',
      `total_number` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '订单总笔数',
      `total_amount` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '订单总金额',
      `status` tinyint UNSIGNED NOT NULL DEFAULT 0 COMMENT '批量状态  1 初始化，2-成功，3-失败, 4-处理中',
      `merchant_id` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '商户id',
      `create_time` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
      `finish_time` int UNSIGNED NULL DEFAULT 0 COMMENT '完成时间',
      `file_content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '文件内容json',
      `generate_all` tinyint UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否已全部生成订单 1-完成，2-未完成',
      PRIMARY KEY (`id`) USING BTREE,
      UNIQUE INDEX `batch_no_unique`(`batch_no`) USING BTREE
    ) ENGINE = InnoDB AUTO_INCREMENT = 22 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

    CREATE TABLE `upload_file_log`  (
      `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
      `file_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '文件名称',
      `account_id` int NULL DEFAULT NULL COMMENT '账号id',
      `account_type` tinyint UNSIGNED NULL DEFAULT NULL COMMENT '账号类型 1-商户',
      `create_time` int UNSIGNED NULL DEFAULT NULL COMMENT '创建时间',
      PRIMARY KEY (`id`) USING BTREE
    ) ENGINE = InnoDB AUTO_INCREMENT = 40 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;

-- 批量付款文件名称
INSERT INTO `global_config`(`config_key`, `config_value`, `remark`, `create_time`, `is_change`) VALUES ('base_batch_transfer_file_name', 'misc/base_batch_transfer.xlsx', '批量付款文件名称', 1622082114, 2);

--
ALTER TABLE `transfer_order`
ADD COLUMN `batch_no` varchar(255) NULL DEFAULT '' COMMENT '批量付款批次号' AFTER `merchant_single_fee`,
ADD COLUMN `batch_row_no` varchar(16) NULL DEFAULT '' COMMENT '批量付款批次行号' AFTER `batch_no`;

-- 商户提现订单表删除字段
ALTER TABLE `merchant_withdraw_order`
DROP COLUMN `merchant_bank_card_id`;