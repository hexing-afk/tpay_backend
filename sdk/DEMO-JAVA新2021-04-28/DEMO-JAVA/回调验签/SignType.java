package tf56.enterprise.demo.common;

/**
 * 验签支持类型
 */
public enum SignType {
    MD5("MD5"),
    RSA("RSA"),
    RSA2("RSA2");
    private String code;

    SignType(String code) {
        this.code = code;
    }

    public String getCode() {
        return code;
    }

    public static SignType valueOfCode(String code) {
        SignType[] values = SignType.values();
        for (SignType value : values) {
            if (value.getCode().equalsIgnoreCase(code)) {
                return value;
            }
        }
        return null;
    }
}
