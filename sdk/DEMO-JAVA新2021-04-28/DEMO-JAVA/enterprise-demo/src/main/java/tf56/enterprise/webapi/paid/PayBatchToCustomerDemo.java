/**
 * Copyright © 2014-2017 TransfarPay.All Rights Reserved.
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
 * 方法说明：商户批量代付到银行卡api；
 * createUser:xufengchun
 * updateDate:20180504
 * updateuser:20501
 * 说明(测试出现余额不足的情况，请通知传化技术人员添加测试账号金额（10W）)
 * 开发人员请仔细阅读文档【每条代付信息共15个字符,不足字符使用^代替,多条字符以|隔开，结尾必须与其他字符隔开】
 */



public class PayBatchToCustomerDemo {

	private final static Logger logger = LoggerFactory.getLogger(PayBatchToCustomerDemo.class);

	private final static String charset = "UTF-8";

	
	private void params1(Map<String, Object> paramsMap) {
		//单条15位(不足用^代替)
		String paydetail = 
				"ONcw00021811264q53510003^"	//流水号在商户端不可重复1^
				+ "分润代付^"						//商品名称2^
				+ "6214835894246071^"			//收款方银行卡号3^
				+ "白芳芳^"						//收款账号姓名4^
				+ "0.1^"						//付款金额5
				+ "个人^"							//银行卡类型6^
				+ "储蓄卡^"						//银行卡借贷类型7^
				+ "招商银行^"						//银行名称8^
				+ "招商银行某支行^"					//支行9^
				+ "浙江省^"						//省10^
				+ "杭州市^"						//市11^
				+ "11111^"						//联行号12^
				+ "描述^"							//描述13^
				+ "备注说明^"						//备注说明14^
				+ "13568475493^"				//手机号码15^
				//手机号码15
				+ "|ONc1020w8WW247q223070004^"//多条信息用|隔开
				+ "分润代付^"
				+ "6217001540013424112^"
				+ "韩思念^"
				+ "0.1^"
				+ "个人^"
				+ "储蓄卡^"
				+ "中国建设银行杭州宝石支行留下分理处^"
				+ "招商银行某某支行^"
				+ "浙江省^"
				+ "杭州市^"
				+ "11111^"
				+ "描述^"
				+ "备注说明^"
				+ "13568475493^"
				;
		//生产测试地址一致
		paramsMap.put("service_id", "tf56enterprise.batchPay.payApply");
		paramsMap.put("tf_timestamp", PaidEnity.getTfTimestamp());
	    paramsMap.put("paydetail", paydetail);
		paramsMap.put("paydate", "20180926");//代付日期(不能小于当天)
		paramsMap.put("batchno", PaidEnity.getDogAk()+PaidEnity.getTfTimestamp());//批量代付批次编号
		paramsMap.put("batchnum", "2");//批量代付总数
		paramsMap.put("batchamount", "0.2");//批量代付总金额
		paramsMap.put("backurl", "http:www.baidu.com/api/notify");
		paramsMap.put("appid", PaidEnity.getAppid());
		paramsMap.put("enterprisecode", PaidEnity.getAccountid());
		paramsMap.put("accountnumber", PaidEnity.getPayaccount());
	}

	public void batchpayToCustomer() {
		//测试环境
		String url = "http://openapitest.tf56.com/service/api";
		//生产环境
		//String url = "https://openapi.tf56.com/service/api";
		Map<String, Object> paramsMap = new HashMap<>();
		params1(paramsMap);
		try {
			String pfxFileName = "D:\\RSAcertificate\\2013001-rafile3ceshi.pfx";// enterprisetest.pfx

			String keyPassword = "741963";
			String originalMessage = ParamUtil.sortMapByKey(paramsMap);
			System.out.println(originalMessage);
			byte[] signMessage = RSAUtils.signMessage(pfxFileName, keyPassword, originalMessage.getBytes());
			String signStr = new String(signMessage);
			paramsMap.put("tf_sign", signStr);
			System.out.println(signStr);
			String result = HttpClient.sendHttpPost(url, paramsMap);
//			String request = HttpClient4Utils.sendHttpRequest(url, map, charset, true);
			System.out.println("request:" + result);
		} catch (Exception e) {
			e.printStackTrace();
		}
	}

	public static void main(String[] args) {
		PayBatchToCustomerDemo batch=new PayBatchToCustomerDemo();
		batch.batchpayToCustomer();
	}

}
