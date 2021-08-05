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
 * 说明：【收银台】查询商户交易结果API 
 * 注意：收银台的订单只有收银台能够查到，代付查询无法获取
 * createUser:20501 
 * updateUser:20501 
 * updateDate:20180507
 */
public class PaySelectResuftDemo2 {

	private final static Logger logger = LoggerFactory.getLogger(PaySelectResuftDemo2.class);

	private final static String charset = "UTF-8";

	private static void params1(Map<String, Object> map) {
		map.put("service_id", "tf56pay.cashier.orderQuery");//网关域名
		map.put("appid", PaidEnity.getAppid());//appid
		map.put("tf_timestamp", PaidEnity.getTfTimestamp());//当前时间戳
		map.put("terminal", "PC");
		map.put("businessnumber", "TESTC20180629142231773236515");// 被查询的单号
		
	}

	public static void main(String[] args) {
		// 测试环境(HTTPS格式就在http后面加s)
		 String url = "http://openapitest.tf56.com/service/api";
		// 生产环境(HTTPS格式就在http后面加s)
		//String url = "https://openapi.tf56.com/service/api";
		Map<String, Object> map = new HashMap<>();
		params1(map);
		try {
//			String pfxFileName = "D:\\Git\\pro-1331001RSA.pfx";// enterprisetest.pfx
			String pfxFileName = "D:\\工作文档\\商户系统文档\\957.pfx";// enterprisetest.pfx

			String keyPassword = "123456";
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
