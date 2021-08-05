/**
 * 
 */
package tf56.enterprise.webapi.paid;

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
 * @date 2018年9月6日
 * @version 1.0 
 * 商户代付_下载商户对账单
 * 
 */
public class PaidDownPayBill {
	private final static Logger logger = LoggerFactory.getLogger(PayForCustomerDemo.class);

	private final static String charset = "UTF-8";

	private static void params1(Map<String, Object> map) {
		// 生产
		// map.put("service_id", "tf56enterprise.enterprise.payForCustomer");
		// 测试
		map.put("service_id", "tf56enterprise.report.downLoadAccountStatement");
		map.put("appid", PaidEnity.getAppid());
		map.put("dog_ak", PaidEnity.getDogAk());
		map.put("tf_timestamp", PaidEnity.getTfTimestamp());
		map.put("enterprisecode", PaidEnity.getAccountid());
		map.put("statementdate", "20180702");//交易日期(不要超过当前日期)
	}

	public static void main(String[] args) {
		// 【生产环境】
		// String url = "http://openapi.tf56.com/service/api";
		// 【测试环境】
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
