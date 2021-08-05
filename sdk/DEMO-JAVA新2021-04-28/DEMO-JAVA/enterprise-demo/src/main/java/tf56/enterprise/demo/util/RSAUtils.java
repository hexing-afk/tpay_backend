/**
 * Copyright © 2014-2017 TransfarPay.All Rights Reserved.
 */
package tf56.enterprise.demo.util;

import com.itrus.cryptorole.CryptoException;
import com.itrus.cryptorole.NotSupportException;
import com.itrus.cryptorole.SignatureVerifyException;
import com.itrus.cryptorole.bc.RecipientBcImpl;
import com.itrus.cryptorole.bc.SenderBcImpl;
import com.itrus.cvm.CVM;
import com.itrus.svm.SignerAndEncryptedDigest;
import org.apache.commons.codec.binary.Base64;

import java.security.Principal;
import java.security.cert.X509Certificate;

/**
 * 描述说明
 * 
 * @version V1.0
 * @author huzz
 * @Date 2017年7月15日 下午4:04:34
 * @since JDK 1.7
 */
public class RSAUtils {

	public static final String CVM_PATH = "src\\resource\\config\\cvm.xml";

	public static final String CER_PATH = "src\\resource\\cafiles\\testca.cer";

	public static byte[] signMessage(String pfxFileName, String keyPassword, byte[] originalMessage) {
		try {
			// 签名流程***********************************************
			// 签名类com.itrus.cryptorole.bc.SenderBcImpl
			SenderBcImpl send = new SenderBcImpl();
			// 使用流程先进行密钥初始化
			send.initCertWithKey(pfxFileName, keyPassword);// pfxFileName PFX证书路径 keyPassword证书密码 。无返回值方法
			// 初始后调用签名方法
			byte[] signMsg = send.signMessage(originalMessage); // originalMessage要签名的原文。
			// 返回签名结果byte[]类型，返回的byte[]类型数据需用再进行编码。
			return Base64.encodeBase64(signMsg);
		} catch (Exception e) {
			System.out.println(e);
			return null;
		}
	}

	public static boolean verifySignature(byte[] originalMessage, byte[] signedData) {

		RecipientBcImpl recipient = new RecipientBcImpl();
        try {
			recipient.initCertWithKey("D:\\RSAcertificate\\2013001-rafile3ceshi.pfx", "741963");
		} catch (NotSupportException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		} catch (CryptoException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
        SignerAndEncryptedDigest signerAndEncryptedDigest = null;
        try {
            signerAndEncryptedDigest = recipient.verifyAndParsePkcs7(originalMessage, Base64.decodeBase64(signedData));
            com.itrus.cert.X509Certificate certificate = com.itrus.cert.X509Certificate.getInstance(signerAndEncryptedDigest.getSigner());
            System.out.println("证书颁发者："+certificate.getSubjectDNString());
            return true;
        } catch (Exception ignored) {
            return false;
        }
		
	}

	public static void main(String[] args) {
		String pfxFileName = "D:\\RSAcertificate\\2013001-rafile3ceshi.pfx";
		String keyPassword = "741963";
		String originalMessage = "C2B8D9F72B7C44D9E6Cb5231A49D92D2";
		byte[] signMessage = signMessage(pfxFileName, keyPassword, originalMessage.getBytes());
		System.out.println("********加签值***************：" + signMessage.toString());
		System.out.println("********验签返回值***************：" + verifySignature(originalMessage.getBytes(), signMessage));
	}
}
