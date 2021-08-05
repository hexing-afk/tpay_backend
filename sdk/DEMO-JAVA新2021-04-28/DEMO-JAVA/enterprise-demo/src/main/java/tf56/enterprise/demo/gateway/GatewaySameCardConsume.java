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
 * 网关前置——同卡消费
 */
public class GatewaySameCardConsume {

	private final static Logger logger = LoggerFactory.getLogger(GatewaySameCardConsume.class);

	private final static String charset = "UTF-8";

	private static void params1(Map<String, Object> map) {
		map.put("appid", PaidEnity.getAppid());
		map.put("tf_timestamp", PaidEnity.getTfTimestamp());
		map.put("service_id", "tf56pay.gateway.sameCardConsume");
		map.put("sign_type", PaidEnity.getSigntype());
		map.put("terminal", "PC");
		map.put("subject", "金融支付");
		map.put("businesstype", "停车费");
		map.put("kind", "电子小票");
		map.put("businessnumber", PaidEnity.getDogAk()+PaidEnity.getTfTimestamp()+PaidEnity.getPayuuid());
		map.put("transactionamount", "0.1");
		map.put("toaccountnumber", PaidEnity.getPayaccount());
		map.put("certcode", "12412312");
		map.put("originalrecordnumber", "hfuisdhgjhsdjghfsjldq");
		map.put("aggregationcode", PaidEnity.getAccountid());
		map.put("clientip", "192.168.0.1");
		map.put("merchantuserid", "10");
		map.put("merchtdevicename", "10010001111");
		map.put("merchtdevicevalue", "10010001111");

	}

	public static void main(String[] args) {
		// 测试环境(HTTPS格式就在http后面加s)
		String url = "https://openapitest.tf56.com/service/api";
		// 生产环境(HTTPS格式就在http后面加s)
		// String url = "https://openapi.tf56.com/service/api";
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
