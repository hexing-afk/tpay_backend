/**
 * Copyright © 2014-2017 TransfarPay.All Rights Reserved.
 */
package tf56.enterprise.demo.gateway;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import tf56.enterprise.demo.util.HttpClient;
import tf56.enterprise.demo.util.HttpClient4Utils;
import tf56.enterprise.demo.util.ParamUtil;
import tf56.enterprise.demo.util.RSAUtils;
import tf56.enterprise.enity.PaidEnity;

import java.util.HashMap;
import java.util.Map;

/**
 * 说明：网关前置_网银支付API 
 * createUser:20501 
 * updateUser:20501 
 * updateDate:20180507
 */
public class GatewayNetSilverDemo {

	private final static Logger logger = LoggerFactory.getLogger(GatewayNetSilverDemo.class);

	private final static String charset = "UTF-8";


	private static void params1(Map<String, Object> map) {

		map.put("service_id", "tf56pay.gateway.bankPay");
		map.put("appid", "1000293");
		map.put("tf_timestamp", PaidEnity.getTfTimestamp());
		map.put("sign_type", PaidEnity.getSigntype());
		map.put("terminal", "PC");
		map.put("fronturl", "http://notify.com");
		map.put("subject", "笔记本");
		map.put("businesstype", "网关和代付");
		map.put("kind", "购物");
		map.put("businessnumber","xxxxxxxsasx" );// 商户流水号
		map.put("transactionamount", "0.01");
		map.put("toaccountnumber", "8802000130487");
		map.put("bankcode", "ICBC");
		map.put("bankaccounttype", "储蓄卡");
		map.put("accountproperty", "对私");// 注意类型,不能搞错【对公/对私】
		map.put("clientip", "192.168.0.1");
		map.put("merchtonline", "0");
	}

	public static void main(String[] args) {
		// 生产
//		 String url = "https://openapi.tf56.com/service/api";
		// 测试
		String url = "http://openapitest.tf56.com/service/api";
		Map<String, Object> map = new HashMap<>();
		params1(map);
		try {
//			String pfxFileName = "D:\\Git\\pro-1331001RSA.pfx";// enterprisetest.pfx
			String pfxFileName = "D:\\RSAcertificate\\2013001-rafile3ceshi.pfx";// enterprisetest.pfx

			String keyPassword = "741963";
			String originalMessage = ParamUtil.sortMapByKey(map);
			System.out.println(originalMessage);
			byte[] signMessage = RSAUtils.signMessage(pfxFileName, keyPassword, originalMessage.getBytes());
			String signStr = new String(signMessage);
			System.out.println(signStr);
			map.put("tf_sign", signStr);
			String result = HttpClient.sendHttpPost(url, map);
//			String request = HttpClient4Utils.sendHttpRequest(url, map, charset, true);
			System.out.println("request:" + result);
		} catch (Exception e) {
			logger.error(e.getMessage(), e);
		}
	}

}
