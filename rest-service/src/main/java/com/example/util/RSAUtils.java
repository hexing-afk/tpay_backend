package com.example.util;

import com.itrus.cryptorole.CryptoException;
import com.itrus.cryptorole.NotSupportException;
import com.itrus.cryptorole.bc.RecipientBcImpl;
import com.itrus.cryptorole.bc.SenderBcImpl;
import com.itrus.svm.SignerAndEncryptedDigest;
import org.apache.commons.codec.binary.Base64;

import java.security.GeneralSecurityException;
import java.util.Arrays;

/**
 * 描述说明
 *
 * @author huzz
 * @version V1.0
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
            System.out.println("证书颁发者：" + certificate.getSubjectDNString());
            return true;
        } catch (Exception ignored) {
            return false;
        }

    }

    public static boolean verifySignature2(byte[] originalMessage, byte[] signedData, String pfx, String pfxSec) {
        RecipientBcImpl recipient = new RecipientBcImpl();
        try {
            recipient.initCertWithKey(pfx, pfxSec.toCharArray());
        } catch (NotSupportException e) {
            // TODO Auto-generated catch block
            System.out.println("1");
            e.printStackTrace();
        } catch (CryptoException e) {
            System.out.println("2");
            // TODO Auto-generated catch block
            e.printStackTrace();
        }
        SignerAndEncryptedDigest signerAndEncryptedDigest = null;
        try {
            System.out.println("3");
            signerAndEncryptedDigest = recipient.verifyAndParsePkcs7(originalMessage, Base64.decodeBase64(signedData));
            System.out.println("4");
            com.itrus.cert.X509Certificate certificate = com.itrus.cert.X509Certificate.getInstance(signerAndEncryptedDigest.getSigner());
            System.out.println("证书颁发者：" + certificate.getSubjectDNString());
            return true;
        } catch (Exception e) {
            System.out.println("5");
            e.printStackTrace();
            return false;
        }
    }

    public static void main(String[] args) {
        String pfxFileName = "D:\\zibozhengshu.pfx";
        String keyPassword = "baibai88";
        String originalMessage = "updatedate=2021-06-11 22:17:47&transactiontype=消费&transactionnumber=1621061122174740007&transactiondate=2021-06-11 22:17:47&transactionamount=5.00&topartyid=1300000395&toaccountnumber=8800013385833&terminal=PC&subject=普通客户维护收费&status=成功&inputdate=2021-06-11 22:17:16&frompartyid=600996268&fromaccountnumber=990305001001&description=普通客户维护收费&businesstype=停车费&businessrecordnumber=1421061122171660012&businessnumber=7162342095239658079509&billamount=5.00&backurl=https://tpay-api.mangopay-test.com/notify/xpay/pay&appid=3657798&";
//		byte[] signMessage = signMessage(pfxFileName, keyPassword, originalMessage.getBytes());

        byte[] signMessage = "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".getBytes();
//		System.out.println("********加签值***************：" + signMessage.toString());
        System.out.println("********验签返回值***************：" + verifySignature2(originalMessage.getBytes(), signMessage,
                pfxFileName, keyPassword));
    }
}
