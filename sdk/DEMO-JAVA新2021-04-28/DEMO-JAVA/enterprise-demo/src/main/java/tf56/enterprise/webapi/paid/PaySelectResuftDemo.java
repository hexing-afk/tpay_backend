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
 * 说明：【代付】查询商户交易结果API 
 * 收银台支付的订单无法查询（一般返回查无此交易）
 * createUser:20501 
 * updateUser:20501 
 * updateDate:20180507
 */
public class PaySelectResuftDemo {

	private final static Logger logger = LoggerFactory.getLogger(PaySelectResuftDemo.class);

	private final static String charset = "UTF-8";

	private static void params1(Map<String, Object> map) {
		map.put("service_id", "tf56pay.enterprise.queryTradeStatus");//网关域名
		map.put("appid", "1331001");//appid
		map.put("tf_timestamp",  PaidEnity.getTfTimestamp());//当前时间戳
		map.put("businessnumber", "1805311443370100028");// 被查询的单号 1805311443370100028

	}

	public static void main(String[] args) {
		// 测试环境(HTTPS格 式就在http后面加s)
//		String url = "http://openapitest.tf56.com/service/api";
		// 生产环境(HTTPS格式就在http后面加s)
		String url = "http://openapi.tf56.com/service/api";
		Map<String, Object> map = new HashMap<>();
		params1(map);
		try {
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
