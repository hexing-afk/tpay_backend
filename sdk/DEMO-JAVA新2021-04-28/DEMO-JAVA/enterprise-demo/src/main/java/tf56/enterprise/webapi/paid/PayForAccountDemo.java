/**
 * Copyright © 2014-2017 TransfarPay.All Rights Reserved.
 */
package tf56.enterprise.webapi.paid;

import java.util.HashMap;
import java.util.Map;

import org.apache.http.client.HttpClient;
import org.apache.http.impl.client.HttpClients;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import tf56.enterprise.demo.util.*;
import tf56.enterprise.enity.PaidEnity;


/**
 * 说明：商户代付到账户
 * createUser:20501
 * updateUser:20501
 * updateDate:20180507
 */
public class PayForAccountDemo {

	private final static Logger logger = LoggerFactory.getLogger(PayForAccountDemo.class);

	private final static String charset = "UTF-8";

	private static void params1(Map<String, Object> map) {
	        map.put("service_id", "tf56enterprise.singlePay.payAccountApply");
	        map.put("appid", PaidEnity.getAppid());
	        map.put("sign_type", PaidEnity.getSigntype());
	        map.put("tf_timestamp", PaidEnity.getTfTimestamp());
	        map.put("businessnumber", PaidEnity.getDogAk()+PaidEnity.getTfTimestamp());
	        map.put("subject", "代发工资");
	        map.put("transactionamount", "1.00");
	        map.put("fromaccountnumber",PaidEnity.getPayaccount());//付款方支付账号
	        map.put("realname", "收款方名称");
	        map.put("toaccountnumber", PaidEnity.getPayaccount());//收款方支付账号
	}

	public static void main(String[] args) {
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
//			String result = HttpClient.sendHttpPost(url, map);
//			String request = HttpClient4Utils.sendHttpRequest(url, map, charset, true);
//			System.out.println("request:" + result);
		} catch (Exception e) {
			logger.error(e.getMessage(), e);
		}
	}

}
