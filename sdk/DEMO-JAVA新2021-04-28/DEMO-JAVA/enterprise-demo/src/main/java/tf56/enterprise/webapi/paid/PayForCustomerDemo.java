/**
 * Copyright © 2014-2017 TransfarPay.All Rights Reserved.
 */
package tf56.enterprise.webapi.paid;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import tf56.enterprise.demo.util.DateUtils;
import tf56.enterprise.demo.util.HttpClient;
import tf56.enterprise.demo.util.HttpClient4Utils;
import tf56.enterprise.demo.util.ParamUtil;
import tf56.enterprise.demo.util.RSAUtils;
import tf56.enterprise.enity.PaidEnity;

import java.util.HashMap;
import java.util.Map;

/**
 * 说明：商户代付到银行卡API 
 * createUser:XNFENGCHUN 
 * update:20501 
 * updateDate:20180504
 */
public class PayForCustomerDemo {

	private final static Logger logger = LoggerFactory.getLogger(PayForCustomerDemo.class);

	private final static String charset = "UTF-8";

	private static void params1(Map<String, Object> map) {
		//生产
//		map.put("service_id", "tf56enterprise.enterprise.payForCustomer");
		//测试
		map.put("service_id", "tf56pay.enterprise.payForCustomer");
		map.put("appid", "1000293");
		map.put("tf_timestamp", DateUtils.strDate("yyyyMMddHHmmss"));
		map.put("businessnumber", "02xxsx01127sssx30121");//随机号码，在商户端不要重复
		map.put("subject", "BILL");
		map.put("transactionamount", "0.01");
		map.put("bankcardnumber", "60131212826515485");
		map.put("bankcardname", "蔡鑫");
		map.put("sign_type", "RSA");
		map.put("bankname", "工商银行");
		map.put("bankcardtype", "个人");
		map.put("bankaccounttype", "储蓄卡");
		map.put("fromaccountnumber", "8802000130487");
	}

	public static void main(String[] args) {
		// 【生产环境】
//		String url = "https://openapi.tf56.com/service/api";
		// 【测试环境】     
		String url = "http://openapitest.tf56.com/service/api";
		Map<String, Object> map = new HashMap<>();
		params1(map);
		try {
			//生产
//			String pfxFileName = "D:\\transfar\\enterprise-demo\\src\\resource\\config\\1331001RSA.pfx";// enterprisetest.pfx
			//测试
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
