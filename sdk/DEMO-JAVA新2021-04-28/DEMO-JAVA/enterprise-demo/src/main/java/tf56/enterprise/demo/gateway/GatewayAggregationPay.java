/**
 * 
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
 * @author 杨xx
 * @date 2018年8月21日  
 * @version 1.0
 * 网关前置——聚合支付
 * 20501
 */
public class GatewayAggregationPay {
	
	
	private final static Logger logger = LoggerFactory.getLogger(GatewayAggregationPay.class);

	private final static String charset = "UTF-8";

	public static void main(String[] args) {
		//[测试地址]【注意千万别看成生产地址】
		//String url = "https://openapitestc.tf56.com/service/api";
		
		//[生产地址]【注意千万别看成测试地址】
		String url = "https://openapi.tf56.com/service/api";
		
		Map<String, Object> map = new HashMap<>();
		map.put("appid", PaidEnity.getAppid());
		map.put("tf_timestamp", PaidEnity.getTfTimestamp());
		map.put("service_id", "tf56pay.gateway.multiPay");
		map.put("sign_type", PaidEnity.getSigntype());
		map.put("terminal", "PC");
		map.put("backurl", "www.baidu.com");
		map.put("fronturl", "www.baidu.com");
		map.put("subject", "电子产品");
		map.put("businesstype", "商户结算");
		map.put("kind", "交易资金结算");
		map.put("businessnumber", PaidEnity.getDogAk()+PaidEnity.getTfTimestamp());//商户名称首字母缩写+时间戳加uuid防重复;
		map.put("transactionamount", "0.01");
		map.put("toaccountnumber", PaidEnity.getPayaccount());
		map.put("bankcode", "ICBC");
		map.put("clientip", "119.165.21.172");
		map.put("merchtonline", "0");
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
