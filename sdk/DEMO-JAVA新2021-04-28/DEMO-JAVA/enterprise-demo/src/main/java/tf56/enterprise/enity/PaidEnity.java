/**
 * 
 */
package tf56.enterprise.enity;

import tf56.enterprise.demo.util.DateUtils;

import java.util.UUID;

/**
 * @author 杨xx
 * @date 2018年9月27日
 * @version 1.0 常用参数模型
 */
public class PaidEnity {
	// 商户appid
	private final static String APPID = "1000293";
	// 支付账号
	private final static String PAYACCOUNT = "8802000130487";
	// 商户编号
	private final static String ACCOUNTID = "5688130133018597";
	// 商户支付密钥
	private final static String DOG_SK = "J0V7Q8eR9Qo4cW19Q742";
	// 商户身份密钥
	private final static String DOG_AK = "6F3tT03BV86cl2p9";
	// 交易时间戳
	private final static String TF_TIMESTAMP = DateUtils.strDate("yyyyMMddHHmmss");
	// 验签方式
	private final static String SIGNTYPE = "RSA";
	// 流水号
	private final static String PAYUUID = UUID.randomUUID().toString().replaceAll("&", "/");

	public static String getAppid() {
		return APPID;
	}

	public static String getPayaccount() {
		return PAYACCOUNT;
	}

	public static String getAccountid() {
		return ACCOUNTID;
	}

	public static String getDogSk() {
		return DOG_SK;
	}

	public static String getDogAk() {
		return DOG_AK;
	}

	public static String getTfTimestamp() {
		return TF_TIMESTAMP;
	}

	public static String getSigntype() {
		return SIGNTYPE;
	}

	public static String getPayuuid() {
		return PAYUUID;
	}
public static void main(String[] args) {
	System.out.println(PAYUUID);
}
}
