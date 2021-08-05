package tf56.enterprise.demo.util;

import static java.util.Comparator.reverseOrder;
import java.nio.charset.StandardCharsets;
import java.security.cert.X509Certificate;
import java.util.Map;
import java.util.TreeMap;
import java.util.Map.Entry;
import org.apache.commons.codec.binary.Base64;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import com.itrus.cryptorole.bc.RecipientBcImpl;
import tf56.enterprise.demo.common.JsonResult;
import tf56.enterprise.demo.common.SignType;


/**
* 描述说明
*
* @version V1.0
* @author fangmw@tythin.com
* @Date 2021年4月28日 上午9:57:24
* @since JDK 1.7
*/
public class RSAUtilsReal {
    private static final String TF_SIGN_PARAM = "tf_sign";
    private final static Logger logger = LoggerFactory.getLogger(RSAUtilsReal.class);

    public static JsonResult verify(Map<String, Object> params, String sign) {
        JsonResult jsonResult = new JsonResult();
        jsonResult.setSignMsg(sign);
        jsonResult.setSignType(SignType.RSA.name());
        try {
            params.remove(TF_SIGN_PARAM);
            String oriData = sortMapByKey(params);
            jsonResult.setOrigMsg(oriData);
            RecipientBcImpl recipient = new RecipientBcImpl();
            //验签方法
            X509Certificate x509Certificate = recipient.verifySignature(oriData.getBytes(StandardCharsets.UTF_8), Base64.decodeBase64(sign));
            jsonResult.setResult(0);
            jsonResult.setSubjectDN(x509Certificate.getSubjectDN().getName());
        } catch (Exception e) {
            logger.warn("签名验证失败，原文已遭篡改。", e);
            jsonResult.setResult(-2);
            jsonResult.setMsg("签名验证失败，原文已遭篡改。");
        }
        return jsonResult;
    }
    private static String sortMapByKey(Map<String, Object> map) {
        if (map == null || map.isEmpty()) {
            return null;
        }
        map.keySet().removeIf(key -> null == map.get(key));
        Map<String, Object> sortedMap = new TreeMap<>(reverseOrder());
        sortedMap.putAll(map);
        StringBuilder result = new StringBuilder();
        for (Entry<String, Object> entry : sortedMap.entrySet()) {
            result.append(entry.getKey()).append("=").append(entry.getValue()).append("&");
        }
        return result.toString();
    }
}
