package tf56.enterprise.demo;

import junit.framework.Test;
import junit.framework.TestCase;
import junit.framework.TestSuite;
import tf56.enterprise.demo.util.DateUtils;
import tf56.enterprise.demo.util.HttpClient;
import tf56.enterprise.demo.util.ParamUtil;
import tf56.enterprise.demo.util.RSAUtils;

import java.util.HashMap;
import java.util.Map;

/**
 * Unit test for simple App.
 */
public class AppTest extends TestCase {

	/**
	 * Create the test case
	 *
	 * @param testName
	 *            name of the test case
	 */
	public AppTest(String testName) {
		super(testName);
	}

	/**
	 * @return the suite of tests being tested
	 */
	public static Test suite() {
		return new TestSuite(AppTest.class);
	}

	/**
	 * Rigourous Test :-)
	 */
	public void testCertSign() {
		String url = "http://openapitest.tf56.com/service/api";
		Map<String, Object> map = new HashMap<>();
//		map.put("service_id", "tf56pay.enterprise.queryEnterpriseAccountBanlance");
//		map.put("appid", "1331001");
//		map.put("tf_timestamp", DateUtils.strDate("yyyyMMddHHmmss"));
//		map.put("sign_type", "RSA");
//		map.put("version", "1");
//		map.put("longitude", "2");
//		map.put("latitude", "3");
//		map.put("clientip", "4");
//		map.put("accountnumber", "8800009486275");
		map.put("service_id", "tf56pay.enterprise.payForCustomer");
		map.put("appid", "1000293");
		map.put("tf_timestamp", DateUtils.strDate("yyyyMMddHHmmss"));
		map.put("businessnumber", "02xxsx01127sssx30121");//随机号码，在商户端不要重复
		map.put("subject", "BILL");
		map.put("transactionamount", "0.01");
		map.put("bankcardnumber", "6013826515485");
		map.put("bankcardname", "蔡鑫");
		map.put("sign_type", "RSA");
		map.put("bankname", "工商银行");
		map.put("bankcardtype", "个人");
		map.put("bankaccounttype", "储蓄卡");
		map.put("fromaccountnumber", "8802000130487");

		try {
//			String pfxFileName = "D:\\Git\\pro-1331001RSA.pfx"; // enterprisetest.pfx
//			String keyPassword = "123456";
			String pfxFileName = "D:\\RSAcertificate\\957.pfx";// enterprisetest.pfx
			String keyPassword = "123456";
			String originalMessage = ParamUtil.sortMapByKey(map);
			System.out.println("-----------originalMessage--------------" + originalMessage);
			byte[] signMessage = RSAUtils.signMessage(pfxFileName, keyPassword, originalMessage.getBytes());
			String signStr = new String(signMessage);
			map.put("tf_sign", signStr);
			System.out.println("------------signStr-------------" + map.get("tf_sign") + "");
			String result = HttpClient.sendHttpPost(url, map);
			System.out.print("------" + result);
		} catch (Exception e) {
			e.printStackTrace();
		}
		assertTrue(true);
	}
}
