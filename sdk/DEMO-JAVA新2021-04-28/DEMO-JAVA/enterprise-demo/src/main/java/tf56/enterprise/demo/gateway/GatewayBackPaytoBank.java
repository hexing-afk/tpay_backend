/**
 * 
 */
package tf56.enterprise.demo.gateway;

import java.util.HashMap;
import java.util.Map;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import tf56.enterprise.demo.util.HttpClient;
import tf56.enterprise.demo.util.HttpClient4Utils;
import tf56.enterprise.demo.util.ParamUtil;
import tf56.enterprise.demo.util.RSAUtils;
import tf56.enterprise.enity.PaidEnity;

/**
 * @author 杨xx
 * @date 2018年9月27日
 * @version 1.0 
 * 网关前置_还款代付
 */
public class GatewayBackPaytoBank {
	private final static Logger logger = LoggerFactory.getLogger(GatewayBackPaytoBank.class);

	private final static String charset = "UTF-8";

	private static void params1(Map<String, Object> map) {

		map.put("appid", PaidEnity.getAppid());
		map.put("tf_timestamp", PaidEnity.getTfTimestamp());
		map.put("service_id", "tf56pay.gateway.backPaytoBank");
		map.put("sign_type", PaidEnity.getSigntype());
		map.put("terminal", "PC");
		map.put("businessnumber", PaidEnity.getDogAk()+PaidEnity.getTfTimestamp()+PaidEnity.getPayuuid());//商户名称首字母缩写+时间戳加uuid防重
		map.put("originalrecordnumber", "5sdfsd432xxxxqqqofjgng5");// 被查询订单号码
		map.put("subject","笔记本");
		map.put("fromaccountnumber",PaidEnity.getPayaccount());
		map.put("transactionamount","1000.00");
		map.put("certificatenumber","372922199601237876");//身份证号
		map.put("bankcardnumber","621583803000644021");//银行卡号
		map.put("bankcardname","李国栋");
		map.put("merchantuserid","1212121");
	}

	public static void main(String[] args) {
		// 生产
		// String url = "https://openapi.tf56.com/service/api";
		// 测试
		String url = "https://openapitest.tf56.com/service/api";
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
