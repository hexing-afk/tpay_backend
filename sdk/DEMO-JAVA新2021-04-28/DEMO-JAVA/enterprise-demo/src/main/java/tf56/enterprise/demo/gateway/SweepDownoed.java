/**
 * Copyright © 2014-2017 TransfarPay.All Rights Reserved.
 */
package tf56.enterprise.demo.gateway;

import java.util.HashMap;
import java.util.Map;
import java.util.Random;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import tf56.enterprise.demo.util.*;
import tf56.enterprise.demo.util.DateUtils;
import tf56.enterprise.demo.util.ParamUtil;
import tf56.enterprise.enity.PaidEnity;

/**
 * 说明：扫码支付【扫码/APP支付消费预下单】 
 * createUser:20501 
 * updateUser:20501 
 * updateDate:20180507
 */
public class SweepDownoed {

	private final static Logger logger = LoggerFactory.getLogger(SweepDownoed.class);

	private final static String charset = "UTF-8";

	private static void params1(Map<String, Object> map) {
		map.put("appid", PaidEnity.getAppid());
		map.put("tf_timestamp", PaidEnity.getTfTimestamp());
		map.put("service_id", "tf56pay.cashier.preOrder");
		map.put("sign_type", PaidEnity.getSigntype());
		map.put("terminal", "PC");
		map.put("subject", "笔记本");
		map.put("businessnumber", PaidEnity.getDogAk()+PaidEnity.getTfTimestamp()+PaidEnity.getPayuuid());//商户流水号
		map.put("billamount", "0.01");
		map.put("transactionamount", "0.01");
		map.put("businesstype", "运费");
		map.put("kind", "物流");
		map.put("toaccountnumber", PaidEnity.getPayaccount());
		map.put("clientip", "119.165.201.172");

	}

	public static void main(String[] args) {
		// 测试环境(HTTPS格式就在http后面加s)
		String url = "http://openapitest.tf56.com/service/api";
		// 生产环境(HTTPS格式就在http后面加s)
		// String url = "http://openapi.tf56.com/service/";
		Map<String, Object> map = new HashMap<>();
		params1(map);
		try {
			map.put("tf_sign", ParamUtil.map2MD5(map));
			System.out.println("----" + map.get("tf_sign") + "");
			map.remove("dog_sk");
			System.out.println(map.get("tf_sign"));
			// true|post请求&&false|get请求
			String request = HttpClient4Utils.sendHttpRequest(url, map, charset, true);
			System.out.println("request:" + request);
		} catch (Exception e) {
			logger.error(e.getMessage(), e);
		}
		
	}

}
