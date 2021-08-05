/**
 * Copyright © 2014-2017 TransfarPay.All Rights Reserved.
 */
package tf56.enterprise.webapi.paid;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import tf56.enterprise.demo.util.MD5;

import java.util.HashMap;
import java.util.Map;

/**
 * 说明：代付业务异步回调商户端验签
 *  createUser:20501 
 *  updateUser:20501 
 *  updateDate:20180507
 *  特别说明【商户收到的回调内容全部放到map里参与加签,不要漏传或误传,尤其是不要粗心的写些空格,或把自己的商户参数配错】
 */
public class PayAttestaTion {

	private final static Logger logger = LoggerFactory.getLogger(PayAttestaTion.class);

	private final static String charset = "UTF-8";

	private static void params1(Map<String, String> map) {
		/*map.put("fromaccountnumber", "990305001001");
		//remark必须注意内容，可能有多个参数存在
		map.put("remark","授予,frompartyid:null,fromaccountnumber:990305001001");
		map.put("status","成功");
		map.put("subject","商户充值");
		map.put("inputdate","2018-08-06 10:23:10");
		map.put("businesstype","充值及回提");
		map.put("transactiontype", "消费");
		map.put("appid","1796001");
		map.put("transactionamount","50000.00");
		map.put("businessnumber", "CZ1533522185674");
		map.put("topartyid","10186311");
		map.put("billamount","50000.00");
		map.put("transactiondate", "2018-08-06 18:23:55");
		map.put("terminal", "PC");
		map.put("fronturl","https://gateway.838pay.com/netPayNotify/tf56PayReturn");
		map.put("toaccountnumber","880001030769"); 
		map.put("transactionnumber","1618080618235550001");
		map.put("businessrecordnumber","1418080610231040002"); 
		map.put("updatedate","2018-08-06 18:23:55"); 
		map.put("backurl","https://gateway.838pay.com/netPayNotify/tf56PayNotify");*/
		
		map.put("timestamp","20181011182000");
		map.put("batchno","B201810111PcdUUg"); 
		map.put("enterprisecode","6688521161015054");
		map.put("sign_type","MD5"); 
		map.put("accountnumber","8800010249355"); 
		map.put("successdetails","B1539253073729791^6212261718003688112^马壮^1.00^成功^Bill->status:成功,remark:null^2918101118194150004^20181011181942");
		map.put("faildetails","B1539253073658187^6212261718003688112^刘杰^11.00^失败^Bill->status:失败,remark:持卡人身份信息验证失败^2918101118194140003^20181011181946"); 
	}

	public static void main(String[] args) {
		Map<String, String> map = new HashMap<>();
		params1(map);
		try {
			//									tf_sign								dog_sk
			boolean verify = MD5.verify(map, "7A66A0C97FAFD8C2C480BC3D440D2D6A", "X32eD09beT5A799p00r3", charset);
			if (verify) {
				System.out.println("验签成功");
			} else {
				System.out.println("验签失败");
			}
		} catch (Exception e) {
			logger.error(e.getMessage(), e);
		}
	}

}
