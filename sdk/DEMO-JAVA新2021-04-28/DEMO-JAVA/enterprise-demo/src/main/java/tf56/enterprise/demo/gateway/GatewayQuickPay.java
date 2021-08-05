/**
 * Copyright © 2014-2017 TransfarPay.All Rights Reserved.
 */
package tf56.enterprise.demo.gateway;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import tf56.enterprise.demo.util.HttpClient;
import tf56.enterprise.demo.util.ParamUtil;
import tf56.enterprise.demo.util.RSAUtils;
import tf56.enterprise.enity.PaidEnity;

import java.util.HashMap;
import java.util.Map;

/**
 * 说明：网关前置_快捷支付
 * createUser:20501 
 * updateUser:20501 
 * updateDate:20180507
 */
public class GatewayQuickPay {

	private final static Logger logger = LoggerFactory.getLogger(GatewayQuickPay.class);

	private final static String charset = "UTF-8";

	private static void params1(Map<String, Object> map) {
		map.put("tf_timestamp", PaidEnity.getTfTimestamp());
		map.put("subject", "笔记本");
		map.put("kind", "实物_数码家电");//消费场景
		map.put("terminal", "PC");
		map.put("transactionamount", "100.00");
		map.put("businesstype", "商家消费");
		map.put("toaccountnumber", PaidEnity.getPayaccount());
		map.put("certcode", "18ed26125cae41e9bd4618cd509ad66d");//客户签约编号
		map.put("appid", PaidEnity.getAppid());
		map.put("service_id","tf56pay.gateway.quickPay");
		map.put("businessnumber", PaidEnity.getDogAk()+PaidEnity.getTfTimestamp()+PaidEnity.getPayuuid());//商户AK+时间戳加uuid防重
		map.put("clientip", "119.165.201.172");
		map.put("sign_type", PaidEnity.getSigntype());
		map.put("merchantuserid", "10010001111");//客户id
	
	}

	public static void main(String[] args) {
		//测试环境(HTTPS格式就在http后面加s)
		String url = "https://openapitest.tf56.com/service/api";
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
