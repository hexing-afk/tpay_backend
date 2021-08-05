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
 * 说明：分页查询交易记录接口API 
 * createUser:20501 
 * updateUser:20501 
 * updateDate:20180507
 */
public class PaySelectPageDemo {

	private final static Logger logger = LoggerFactory.getLogger(PaySelectPageDemo.class);

	private final static String charset = "UTF-8";

	private static void params1(Map<String, Object> map) {
		map.put("service_id", "tf56pay.enterprise.queryTransactionDetails");
		map.put("appid", PaidEnity.getAppid());
		map.put("tf_timestamp", PaidEnity.getTfTimestamp());
		map.put("startdate", "2018-07-30 00:00:00");
		map.put("enddate", "2018-08-01 10:00:00");
		map.put("pagesize", "10");
		map.put("skipcount", "0");
	}

	public static void main(String[] args) {
		// 测试环境
		String url = "http://openapitest.tf56.com/service/api";
		// 生产环境
		// String url = "http://openapi.tf56.com/service/api";
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
