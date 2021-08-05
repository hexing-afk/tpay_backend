/**
 * Copyright © 2014-2017 TransfarPay.All Rights Reserved.
 */
package tf56.enterprise.webapi.paid;

import java.util.HashMap;
import java.util.Map;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import tf56.enterprise.demo.util.*;
import tf56.enterprise.enity.PaidEnity;


/**
 * 说明：商户批量代付到账户
 * createUser:20501
 * updateUser:20501
 * updateDate:20180507
 */
public class PayBatchAccountDemo {

	private final static Logger logger = LoggerFactory.getLogger(PayBatchAccountDemo.class);

	private final static String charset = "UTF-8";

	private static void params1(Map<String, Object> map) {
		String paydetail="122W567890^"			//商户流水号
							+ "8802000087147^"	//收款方账户
							+ "收款名称^"			//收款名称
							+ "10000.00^"		//付款金额
							+ "xx^"				//描述
							+"|12345678E0^"		//商户流水号
							+ "8802000087147^"	//收款方账户
							+ "收款名称^"			//收款名称
							+ "10000.00^"		//付款金额
							+ "xx^";				//描述
	        map.put("service_id", "tf56enterprise.batchPay.payAccountApply");
	        map.put("appid", PaidEnity.getAppid());
	        map.put("tf_timestamp", PaidEnity.getTfTimestamp());
	        map.put("terminal", "PC");
	        map.put("enterprisecode", PaidEnity.getAccountid());
	        map.put("accountnumber", PaidEnity.getPayaccount());//收款方支付账号
	        map.put("batchno", "2018E80W4q2q2");//批量代付批次编号
	        map.put("batchnum", "2");//批量笔数
	        map.put("batchamount", "20000.00");//总金额
	        map.put("paydate", "20180928");//代付日期(不能小于当天)
	        map.put("paydetail", paydetail);
	}

	public static void main(String[] args) {
		//测试环境地址
		String url = "http://openapitest.tf56.com/service/api";
		// 【生产环境】
		//String url = "http://openapi.tf56.com/service/api";
		Map<String, Object> map = new HashMap<>();
		params1(map);
		try {
			
			String pfxFileName = "D:\\RSAcertificate\\2013001-rafile3ceshi.pfx";// enterprisetest.pfx

			String keyPassword = "741963";
			String originalMessage = ParamUtil.sortMapByKey(map);
			System.out.println(originalMessage);
			byte[] signMessage = RSAUtils.signMessage(pfxFileName, keyPassword, originalMessage.getBytes());
			String signStr = new String(signMessage);
			map.put("tf_sign", signStr);
			System.out.println(signStr);
			String result = HttpClient.sendHttpPost(url, map);
//			String request = HttpClient4Utils.sendHttpRequest(url, map, charset, true);
			System.out.println("request:" + result);
		} catch (Exception e) {
			logger.error(e.getMessage(), e);
		}
	}

}
