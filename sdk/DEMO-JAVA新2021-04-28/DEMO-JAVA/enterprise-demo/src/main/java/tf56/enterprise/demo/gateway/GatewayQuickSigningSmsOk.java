/**
 * Copyright © 2014-2017 TransfarPay.All Rights Reserved.
 */
package tf56.enterprise.demo.gateway;

import java.util.HashMap;
import java.util.Map;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import tf56.enterprise.demo.util.*;
import tf56.enterprise.enity.PaidEnity;

/**
 * 说明：网关前置_快捷签约确认
 * createUser:20501 
 * updateUser:20501 
 * updateDate:20180507
 */
public class GatewayQuickSigningSmsOk {

	private final static Logger logger = LoggerFactory.getLogger(GatewayQuickSigningSmsOk.class);

	private final static String charset = "UTF-8";

	private static void params1(Map<String, Object> map) {
		map.put("appid", PaidEnity.getAppid());
		map.put("tf_timestamp", PaidEnity.getTfTimestamp());
		map.put("service_id","tf56pay.gateway.quickPayConfirm");
		map.put("sign_type", PaidEnity.getSigntype());
		map.put("terminal", "PC");
		map.put("verifycode", "794854");
		map.put("clientip", "255.255.10.1");
		map.put("businessrecordnumber", "支付订单号");
	}

	public static void main(String[] args) {
		//测试环境
		String url = "https://openapitest.tf56.com/service/api";
		//生产环境
		//String url = "https://openapi.tf56.com/service/api";
		Map<String, Object> map = new HashMap<>();
		params1(map);
		try {
			String pfxFileName = "D:\\work\\a\\2013001-rafile3ceshi.pfx";// enterprisetest.pfx

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
