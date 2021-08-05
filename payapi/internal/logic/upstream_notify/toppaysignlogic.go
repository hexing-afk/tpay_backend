package upstream_notify

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/core/logx"
	"tpay_backend/payapi/internal/svc"
)

type ToppaySignLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewToppaySignLogic(ctx context.Context, svcCtx *svc.ServiceContext) ToppaySignLogic {
	return ToppaySignLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ToppaySignLogic) ToppaySign(body []byte) error {
	//真实示例：
	//businessrecordnumber=3121061116310940024\u0026certcode=52412f935e7f48edb7e61cc03e0dcb2c\u0026backurl=https%3A%2F%2Ftpay-api.mangopay-test.com%2Fnotify%2Fxpay%2Fpay\u0026appid=3657798\u0026businessnumber=7162340018716291053256\u0026signtype=%E7%AD%BE%E7%BA%A6\u0026bankname=%E4%B8%AD%E5%9B%BD%E9%93%B6%E8%A1%8C\u0026tf_sign=MIIF9gYJKoZIhvcNAQcCoIIF5zCCBeMCAQExCzAJBgUrDgMCGgUAMAsGCSqGSIb3DQEHAaCCBA4wggQKMIIC8qADAgECAhRBE9W7jPwhBH3J4h%2F71BPjt8fFhTANBgkqhkiG9w0BAQsFADBwMQswCQYDVQQGEwJDTjEhMB8GA1UECgwY5Lyg5YyW5pSv5LuY5pyJ6ZmQ5YWs5Y%2B4MRgwFgYDVQQLDA%2FmlK%2Fku5jlvIDlj5Hpg6gxJDAiBgNVBAMMG%2BS8oOWMluaUr%2BS7mOaciemZkOWFrOWPuCBDQTAeFw0xOTEyMDIwOTQ1NTRaFw0yMDEyMDEwOTQ1NTRaMIGHMSIwIAYJKoZIhvcNAQkBFhNqeGwyMDA5MjAwOUAxNjMuY29tMSQwIgYDVQQDDBtSU0E2Njg4MTMwMTMzMDE0ODkzXzEzMzEwMDExGDAWBgNVBAsMD%2BaUr%2BS7mOW8gOWPkemDqDEhMB8GA1UECgwY5Lyg5YyW5pSv5LuY5pyJ6ZmQ5YWs5Y%2B4MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAxRHLhm5jNQadGV9FfE%2BfFoWKe7oq%2BG0spQsDGCL2URh7hhiitDDBMO%2BDLFKCBhLXTKwO%2F1MzvKedXM3kt17gAjB8pfQSmbhA%2Fwi8T9elZdAvfcIPNwIYeitAgOkjL8b6WYJrvuhXrMHThXsENrx5kD5lXrUwRDtJ%2BABxU2oGJNoJbCXFIZTYsRRJBdbRDIEDR44qwVfPG3sb2vpmtaKua7CaqJL4AcR2SbhrA2Tja00oAWsXWCtmBOZMWgmEmIMYH42az2SH0k%2F8ALXqoLHyWq7MvkEyqd1Pi2vngdcxp0m%2F1s2548Lu7BDY%2FeHpIdk5QpODJbzTzJXk%2BtkC%2FLAO1QIDAQABo4GDMIGAMAkGA1UdEwQCMAAwCwYDVR0PBAQDAgWgMGYGA1UdHwRfMF0wW6BZoFeGVWh0dHA6Ly90b3BjYS5pdHJ1cy5jb20uY24vcHVibGljL2l0cnVzY3JsP0NBPTdGN0MxRDIxRDQ0RUJERDFGQkI4RDM5QzcwRDE1RTcwMkNBNkJCNzEwDQYJKoZIhvcNAQELBQADggEBALHP1nHHeM27ZOofeCuo5FsI1F1RVF9%2BgVMgX056Zj2jnPXNREbtq1IzP%2BHpzDkTSJQepS5TWdsbFwC4ghMmeg3KYDOsRXVOZVrqGRXg%2FEE0XDRiJXvTzcVLpg5NTQIJwwADJwuWHXapPnUNk6TEY701wXpMAN6YYU%2BAYnuhm1BZRQZjZVb79na%2F34%2BkdKluDoIafa43QkJNScH1TD7hkbYbdH4MhubGLSgIERpHW%2FXOcxZSREkyE%2FaEEmN6odthMJFf4zc%2B6UOyzYE56blDvsbeJxZnTaGfH0OQp0Pz%2BoqnuxCeRS%2BaYNMYxLrPycQJjWxtQa2KZf67LnNExE8KztAxggGwMIIBrAIBATCBiDBwMQswCQYDVQQGEwJDTjEhMB8GA1UECgwY5Lyg5YyW5pSv5LuY5pyJ6ZmQ5YWs5Y%2B4MRgwFgYDVQQLDA%2FmlK%2Fku5jlvIDlj5Hpg6gxJDAiBgNVBAMMG%2BS8oOWMluaUr%2BS7mOaciemZkOWFrOWPuCBDQQIUQRPVu4z8IQR9yeIf%2B9QT47fHxYUwCQYFKw4DAhoFADANBgkqhkiG9w0BAQEFAASCAQCyEWLrX8bQ7CjbJUEKEKL2Y7JciM4PUGE5kQHh%2B1mAadACDn9RwewoXNuqAHi2AKLv%2FeyOMO0fIs%2BEdpC%2BrkMGt%2ByVzZjVFkJst%2BXdbIS9K5ZIP1PIxVVobk7R97j0o6%2BEsMWpu2xdBQpqkD2mSd19CvwK2zDebEJPE%2BpMXL8J6skHMcILlzIh%2FgkYO2sriNqGAiR3TNDhfDAUY%2BNA2aSFw4Kw5GNWQzwxHzuImsxNpNsA8%2BFm69MUebgU3LEA%2BH%2Fsw52HBBnszNgiTsI1VrxlptuASRlTpGkcXtX9zngjboEvXx6socG03j8zy3Xt2cZsx%2Fs86gI6TlfmllAoweV%2F\u0026signdate=2021-06-11+16%3A32%3A35\u0026status=%E6%88%90%E5%8A%9F

	bodyStr := string(body)
	bodyStr = strings.ReplaceAll(bodyStr, "\\u003c", "<")
	bodyStr = strings.ReplaceAll(bodyStr, "\\u003e", ">")
	bodyStr = strings.ReplaceAll(bodyStr, "\\u0026", "&")

	//解析url字符串
	u, err := url.ParseQuery(bodyStr)
	if err != nil {
		panic(err)
	}

	for k, v := range u {
		fmt.Printf("u[%v]=%v\n", k, v)
	}

	//todo 将签约编号存入redis中去
	//1.查询签约订单号对应的卡号
	cardNo, err := l.svcCtx.Redis.Get(context.TODO(), GetBusinessrecordnumberRedisKey(u["businessrecordnumber"])).Result()
	if err != nil {
		logx.Errorf("查询redis 出错,key=%v", GetBusinessrecordnumberRedisKey(u["businessrecordnumber"]))
	}

	//2.设置签约编号
	ret, err := l.svcCtx.Redis.SetEX(context.TODO(), GetRedisKey(cardNo), nil, time.Minute*20).Result()
	logx.Infof("ret=%v,err=%v", ret, err)

	return nil
}

func GetRedisKey(t interface{}) string {
	s := fmt.Sprintf("QPaySign%s", t)
	s = strings.Replace(s, "*", "", -1) // 去掉符号 *
	s = strings.Replace(s, ".", "", -1) // 去掉符号 .
	return s
}

func GetBusinessrecordnumberRedisKey(t interface{}) string {
	s := fmt.Sprintf("businessrecordnumber:%s", t)
	s = strings.Replace(s, "*", "", -1) // 去掉符号 *
	s = strings.Replace(s, ".", "", -1) // 去掉符号 .
	return s
}
