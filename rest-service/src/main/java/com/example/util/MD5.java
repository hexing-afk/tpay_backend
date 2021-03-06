package com.example.util;

import java.io.UnsupportedEncodingException;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.Arrays;
import java.util.HashMap;
import java.util.Map;

/**
 * 功能：md5签名 
 * 修改日期：2016-9-20 10:19:23 
 * 说明： 以下代码只是为了方便商户测试而提供的样例代码，商户可以根据自己网站的需要，按照技术文档编写,并非一定要使用该代码。
 */

public class MD5 {

	public static final String TF_TIME_PARAM = "tf_timestamp";

	private static final String TF_SIGN_PARAM = "tf_sign";

	private static final String TF_DOG_SK = "dog_sk";

	public static void main(String[] args) {
		System.out.println("启动测试MD5例子");
	}

	/**
	 * 签名字符串
	 * 
	 * @param params
	 * @param key
	 * @param charset
	 * @return
	 */
	public static String sign(Map<String, String> params, String key, String charset) {
		// 1. tf_sign不加入签名
		params.remove(TF_SIGN_PARAM);
		// 2. 加入sk
		params.put(TF_DOG_SK, key);
		// 3. 对请求参数的值进行排序
		String keyString = createLinkString(params);
		System.out.println(keyString);
		// 4. 生成md5
		return generateMd5(keyString, charset);
	}
	
	/**
	 * 签名字符串，外部传入时间戳
	 * 
	 * @param params
	 * @param charset
	 * @return
	 */
	public static String signWithTime(Map<String, String> params, String charset) {
		// 1. 对请求参数的值进行排序
		String keyString = createLinkString(params);
		// 2. 生成md5签名
		return generateMd5(keyString, charset);
	}
	

	/**
	 * 生成md5
	 *
	 * @param charset
	 * @return
	 */
	public static String generateMd5(String keyString, String charset) {
		String result = "";
		try {
			MessageDigest md5 = MessageDigest.getInstance("MD5");
			md5.update(keyString.getBytes(charset));
			byte[] temp;
			temp = md5.digest("".getBytes(charset));
			for (int i = 0; i < temp.length; i++) {
				result += Integer.toHexString((0x000000ff & temp[i]) | 0xffffff00).substring(6);
			}
			result = result.toUpperCase();
			System.out.println(result);
		} catch (NoSuchAlgorithmException e) {
			throw new RuntimeException("MD5签名过程中出现错误" + e);
		} catch (UnsupportedEncodingException e) {
			throw new RuntimeException("MD5签名过程中出现错误,指定的编码集不对,您目前指定的编码集是:" + charset);
		}
		return result;
	}

	/**
	 * 除去数组中的空值和签名参数
	 * 
	 * @param sArray
	 *            签名参数组
	 * @return 去掉空值与签名参数后的新签名参数组
	 */
	public static Map<String, String> paraFilter(Map<String, String> sArray) {

		Map<String, String> result = new HashMap<String, String>();

		if (sArray == null || sArray.size() <= 0) {
			return result;
		}

		for (String key : sArray.keySet()) {
			String value = sArray.get(key);
			if (value == null || value.equals("") || key.equalsIgnoreCase("sign")) {
				continue;
			}
			result.put(key, value);
		}

		return result;
	}

	/**
	 * 把数组所有元素排序，并参数值拼接成字符串
	 * 
	 * @param params
	 *            需要排序并参与字符拼接的参数组
	 * @return 拼接后字符串
	 */
	public static String createLinkString(Map<String, String> params) {
		Object[] keys = params.keySet().toArray();
		Arrays.sort(keys);
		StringBuffer sb = new StringBuffer();
		for (int i = keys.length - 1; i >= 0; i--) {
			sb.append(params.get(keys[i]));
		}
		return sb.toString();
	}

	/**
	 * 校验签名
	 * 
	 * @param params
	 * @param sign
	 * @param key
	 * @param charset
	 * @return
	 */
	public static boolean verify(Map<String, String> params, String sign, String key, String charset) {
		// 1. tf_sign不加入签名
		params.remove(TF_SIGN_PARAM);
		// 2. 加入sk
		params.put(TF_DOG_SK, key);
		String keyString = createLinkString(params);
		System.out.println("verify | 拼装结果createLinkString:"+keyString);
		// 3. 生成md5
		String md5 = generateMd5(keyString, charset);
		System.out.println("verify | 生成generateMd5:"+md5);
		if (md5.equals(sign)) {
			return true;
		} else {
			return false;
		}
	}
	
	 
}