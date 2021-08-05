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
 * @version 1.0 网关前置——分账订单查询
 */
public class GatewayShareOrderQuery {
	private final static Logger logger = LoggerFactory.getLogger(GatewayShareOrderQuery.class);

	private final static String charset = "UTF-8";

	private static void params1(Map<String, Object> map) {
		map.put("appid", PaidEnity.getAppid());
		map.put("tf_timestamp", PaidEnity.getTfTimestamp());
		map.put("service_id", "tf56pay.gateway.shareOrderQuery");
		map.put("sign_type", PaidEnity.getSigntype());
		map.put("terminal", "PC");
		map.put("businessnumber", "1418062713044720002");//被查询订单号

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
