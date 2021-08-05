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
 * 说明：网关前置_网银支付_订单查询API 
 * createUser:20501 
 * updateUser:20501 
 * updateDate:20180507
 */
public class GatewayNetSilverSelectDemo {

	private final static Logger logger = LoggerFactory.getLogger(GatewayNetSilverSelectDemo.class);

	private final static String charset = "UTF-8";
	
	private static void params1(Map<String, Object> map) {
		
		map.put("appid", PaidEnity.getAppid());
		map.put("tf_timestamp", PaidEnity.getTfTimestamp());
		map.put("service_id", "tf56pay.gateway.orderQuery");
		map.put("sign_type", PaidEnity.getSigntype());
		map.put("terminal","PC");
		map.put("businessnumber", "5sdfsd432xxxxqqqofjgng5");//被查询订单号码
	}

	public static void main(String[] args) {
		//生产
		String url = "https://openapi.tf56.com/service/api";
		//测试
		//String url = "https://openapitest.tf56.com/service/api";
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
