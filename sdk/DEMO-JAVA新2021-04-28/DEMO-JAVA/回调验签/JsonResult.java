package tf56.enterprise.demo.common;

import java.io.Serializable;

/**
 * 加签验签返回结果<br>
 * result 0-成功，其他失败
 * msg 执行错误信息
 * signType 签名类型，如MD5，ItrusRSA
 * origMsg 加签/验签原文
 * signMsg 加签/验签密文
 * 〈itrusRSAVerify msg〉
 *
 * @author z.z.hu
 * @create 2018/10/16
 * @since 1.0.0
 */
public class JsonResult implements Serializable {

    private static final long serialVersionUID = -2502038518362321067L;
    private String msg;
    private int result;
    private String signType;
    private String origMsg;
    private String signMsg;
    private String subjectDN;

    public JsonResult() {
        this.result = -1;
    }

    public JsonResult(String msg, int result, String signType, String origMsg, String signMsg) {
        this.msg = msg;
        this.result = result;
        this.signType = signType;
        this.origMsg = origMsg;
        this.signMsg = signMsg;
    }

    public String getMsg() {
        return msg;
    }

    public void setMsg(String msg) {
        this.msg = msg;
    }

    public int getResult() {
        return result;
    }

    public void setResult(int result) {
        this.result = result;
    }

    public String getSignType() {
        return signType;
    }

    public void setSignType(String signType) {
        this.signType = signType;
    }

    public String getOrigMsg() {
        return origMsg;
    }

    public void setOrigMsg(String origMsg) {
        this.origMsg = origMsg;
    }

    public String getSignMsg() {
        return signMsg;
    }

    public void setSignMsg(String signMsg) {
        this.signMsg = signMsg;
    }

    public String getSubjectDN() {
        return subjectDN;
    }

    public void setSubjectDN(String subjectDN) {
        this.subjectDN = subjectDN;
    }

    @Override
    public String toString() {
        final StringBuilder sb = new StringBuilder("{");
        sb.append("\"msg\":\"")
                .append(msg).append('\"');
        sb.append(",\"result\":")
                .append(result);
        sb.append(",\"signType\":\"")
                .append(signType).append('\"');
        sb.append(",\"origMsg\":\"")
                .append(origMsg).append('\"');
        sb.append(",\"signMsg\":\"")
                .append(signMsg).append('\"');
        sb.append(",\"subjectDN\":\"")
                .append(subjectDN).append('\"');
        sb.append('}');
        return sb.toString();
    }
}
