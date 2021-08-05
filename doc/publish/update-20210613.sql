CREATE TABLE tpay.bank (
	id int auto_increment NOT NULL,
	bank_name varchar(100) NULL COMMENT '银行名',
	single_amount BIGINT DEFAULT 0 NULL COMMENT '单笔金额',
	CONSTRAINT bank_pk PRIMARY KEY (id)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_0900_ai_ci;

ALTER TABLE tpay.bank ADD CONSTRAINT bank_un UNIQUE KEY (bank_name);

INSERT INTO tpay.bank (bank_name,single_amount) VALUES
	 ('中国工商银行',1000000),
	 ('中国银行',1000000),
	 ('中国邮储银行',800000),
	 ('中信银行',5000000),
	 ('中国光大银行',500000),
	 ('上海银行',1000000),
	 ('浦发银行',1000000),
	 ('华夏银行',1000000),
	 ('平安银行',1000000),
	 ('恒丰银行',5000000);
INSERT INTO tpay.bank (bank_name,single_amount) VALUES
	 ('渤海银行',3000000),
	 ('浙商银行',2000000),
	 ('宁波银行',2000000);

ALTER TABLE tpay.platform_channel ADD start_time BIGINT DEFAULT 0 NULL COMMENT '开始时间';
ALTER TABLE tpay.platform_channel ADD end_time BIGINT DEFAULT 0 NULL COMMENT '结束时间';
ALTER TABLE tpay.platform_channel ADD start_amount BIGINT DEFAULT 0 NULL COMMENT '开始金额';
ALTER TABLE tpay.platform_channel ADD end_amount BIGINT DEFAULT 0 NULL COMMENT '结束金额';


