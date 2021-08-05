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
 * 说明：网关前置_快捷签约【页面跳转】
 * createUser:20501 
 * updateUser:20501 
 * updateDate:20180507
 */
public class GatewayQuickSignhtml {

	private final static Logger logger = LoggerFactory.getLogger(GatewayQuickSignhtml.class);

	private final static String charset = "UTF-8";

	private static void params1(Map<String, Object> map) {
		map.put("appid", PaidEnity.getAppid());
		map.put("tf_timestamp", PaidEnity.getTfTimestamp());
		map.put("service_id","tf56pay.gateway.quickSignPage");
		map.put("sign_type", PaidEnity.getSigntype());
		map.put("terminal", "PC");
		map.put("backurl", "www.back.url");
		map.put("businessnumber", PaidEnity.getDogAk()+PaidEnity.getTfTimestamp()+PaidEnity.getPayuuid());
		map.put("bankcardnumber", "6259063489450034");
		map.put("bankcardname", "中国银行");
		map.put("certificatenumber", "340111199004030017");
		map.put("certtype", "01");
		map.put("bankmobilenumber", "13866137516");
		map.put("clientip", "119.165.201.172");
		map.put("merchantuserid", "10010001111");//客户id
	
	}

	public static void main(String[] args) {
		//测试环境(HTTPS格式就在http后面加s)
		String url = "https://openapitestc.tf56.com/service/api";
		//生产环境(HTTPS格式就在http后面加s)
		//String url = "https://openapi.tf56.com/service/api";
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
