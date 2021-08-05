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
 * 说明：网关前置_快捷签约(无页面跳转)
 * createUser:20501 
 * updateUser:20501 
 * updateDate:20180507
 */
public class GatewayQuickSigning {

	private final static Logger logger = LoggerFactory.getLogger(GatewayQuickSigning.class);

	private final static String charset = "UTF-8";

	private static void params1(Map<String, Object> map) {
		map.put("tf_timestamp", "20210531135619"); //PaidEnity.getTfTimestamp());
		map.put("service_id","tf56pay.gateway.quickSign" );
		map.put("terminal", "PC");
		map.put("version", "01");
		map.put("backurl", "www.baidu.com");
		map.put("businessnumber", "test0002");
		map.put("bankcardnumber", "6212261001005677700");
		map.put("bankcardname","老王");
		map.put("certificatenumber", "431025199906235623");
		map.put("certtype", "01");
		map.put("bankmobilenumber", "13866137516");
		map.put("clientip", "127.0.0.1");
		map.put("merchantuserid", "1861");
		map.put("cvv2", "888");//信用卡必填/非信用卡可不填
		map.put("expiredate", "0323");//信用卡必填/非信用卡可不填
		map.put("merchtdevicename", "28:c6:3f:3d:cc:e7");
		map.put("merchtdevicevalue", "28:c6:3f:3d:cc:e7");
		map.put("appid", "1000293");//
		map.put("dog_sk", "J0V7Q8eR9Qo4cW19Q742");
		map.put("sign_type", "RSA");//RSA
	}

	public static void main(String[] args) {
		//测试环境
//		String url = "https://openapitest.tf56.com/service/api";
		//生产环境
		String url = "";// "https://openapi.tf56.com/service/api";
		Map<String, Object> map = new HashMap<>();
		params1(map);
		try {
			String pfxFileName = "D:\\work\\a\\2013001-rafile3ceshi.pfx";// enterprisetest.pfx

			String keyPassword = "741963";
			String originalMessage = ParamUtil.sortMapByKey(map);
			System.out.println(originalMessage);
//			originalMessage = "appid=1000293&bankcardname=asdfasf&bankcardnumber=65465462184684&bankmobilenumber=1231312314&businessnumber=2021053113435421026723&certificatenumber=2342342425&certtype=01&clientip=127.0.0.1&merchantuserid=6688531146016889&service_id=tf56pay.gateway.quickSign&sign_type=RSA&terminal=Android&tf_timestamp=20210531134354&";
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
