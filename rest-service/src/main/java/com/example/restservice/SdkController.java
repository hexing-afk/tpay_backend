package com.example.restservice;

import com.example.config.SdkConfig;
import com.example.model.*;
import com.example.util.DateUtils;
import com.example.util.HttpClient4Utils;
import com.example.util.ParamUtil;
import com.example.util.RSAUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import java.nio.charset.StandardCharsets;
import java.util.HashMap;
import java.util.Map;
import java.util.TreeMap;

@RestController
public class SdkController {

    @Autowired
    private SdkConfig sdkConfig;

    @GetMapping("/v")
    public String v() {
        return sdkConfig.toString();
    }

    @GetMapping("/v2")
    public String v2(String key) {
        return sdkConfig.getKey(key).toString();
    }

    @PostMapping("/quickSign")
    public SdkResp quickSign(@RequestBody SdkReq body) {
        SdkConfig.KeyConfig s = sdkConfig.getKey(body.getAppid());
        String url = sdkConfig.getGateway();
//        System.out.printf("body=%s\n", body);

        Map<String, Object> map = new HashMap<>();
        map.put("tf_timestamp", DateUtils.strDate("yyyyMMddHHmmss"));
        map.put("service_id", "tf56pay.gateway.quickSign");
        map.put("terminal", "PC");
        map.put("version", "01");
        map.put("appid", body.getAppid());//
        map.put("sign_type", "RSA");//RSA

        map.put("backurl", body.getBackUrl());
        map.put("businessnumber", body.getBusinessnumber());
        map.put("bankcardnumber", body.getBankcardnumber());
        map.put("bankcardname", body.getBankcardname());
        map.put("certificatenumber", body.getCertificatenumber());
        map.put("certtype", "01");
        map.put("bankmobilenumber", body.getBankmobilenumber());
        map.put("clientip", body.getClientip());

        map.put("merchantuserid", s.getMerchantUserid());

//        {"biz_code":"GPBIZ_011001","biz_msg":"没有可用的签约渠道","code":"GP_00","msg":"已收到商户请求，接口正常"}
        String result = "{\"code\":\"GP_00\",\"msg\":\"\",\"biz_code\":\"GPBIZ_00\",\"biz_msg\":\"\",\"data\":{\"sign_type\":\"RSA\",\"tf_sign\":\"1\",\"businessnumber\":\"1\",\"businessrecordnumber\":\"1\",\"certcode\":\"2\",\"status\":\"1\"}}";
        try {
            String pfxFileName = s.getPfxFile();
            String keyPassword = s.getPfxSecret();
            String originalMessage = ParamUtil.sortMapByKey(map);
            System.out.println(originalMessage);
            byte[] signMessage = RSAUtils.signMessage(pfxFileName, keyPassword, originalMessage.getBytes());
            String signStr = new String(signMessage);
//            System.out.println("tf_sign:" + signStr);
            map.put("tf_sign", signStr);
            switch (sdkConfig.getMode()) {
                case "test":
                    break;
                case "normal":
                    result = HttpClient4Utils.sendHttpRequest(url, map, "utf8", true);
                    System.out.println("request:" + result);
                    break;
                default:
                    break;
            }
        } catch (Exception e) {
            e.printStackTrace();
        }

        return new SdkResp(result, "");
    }

    @PostMapping("/quickSignConfirm")
    public SdkResp quickSignConfirm(@RequestBody QuickSignConfirmReq body) {
        SdkConfig.KeyConfig s = sdkConfig.getKey(body.getAppid());
        String url = sdkConfig.getGateway();
//        System.out.printf("body=%s\n", body);

        Map<String, Object> map = new HashMap<>();
        map.put("tf_timestamp", DateUtils.strDate("yyyyMMddHHmmss"));
        map.put("service_id", "tf56pay.gateway.quickSignConfirm");
        map.put("terminal", "PC");
        map.put("version", "01");
        map.put("appid", s.getAppid());//
        map.put("sign_type", "RSA");//RSA

        map.put("businessrecordnumber", body.getBusinessrecordnumber());
        map.put("verifycode", body.getVerifycode());
        map.put("clientip", body.getClientip());

//        String result = "";
        String result = "{\"code\":\"GP_00\",\"msg\":\"\",\"biz_code\":\"GPBIZ_00\",\"biz_msg\":\"\",\"data\":{\"sign_type\":\"RSA\",\"tf_sign\":\"1\",\"businessnumber\":\"1\",\"businessrecordnumber\":\"1\",\"certcode\":\"2\",\"status\":\"1\"}}";
        try {
            String pfxFileName = s.getPfxFile();
            String keyPassword = s.getPfxSecret();
            String originalMessage = ParamUtil.sortMapByKey(map);
            System.out.println(originalMessage);
            byte[] signMessage = RSAUtils.signMessage(pfxFileName, keyPassword, originalMessage.getBytes());
            String signStr = new String(signMessage);
//            System.out.println("tf_sign:" + signStr);
            map.put("tf_sign", signStr);
            switch (sdkConfig.getMode()) {
                case "test":
                    break;
                case "normal":
                    result = HttpClient4Utils.sendHttpRequest(url, map, "utf8", true);
                    System.out.println("request:" + result);
                    break;
                default:
                    break;
            }
        } catch (Exception e) {
            e.printStackTrace();
        }

        return new SdkResp(result, "");
    }

    @PostMapping("/signQuery")
    public SdkResp signQuery(@RequestBody SignQueryReq body) {
        SdkConfig.KeyConfig s = sdkConfig.getKey(body.getAppid());
        String url = sdkConfig.getGateway();
//        System.out.printf("body=%s\n", body);

        Map<String, Object> map = new HashMap<>();
        map.put("tf_timestamp", DateUtils.strDate("yyyyMMddHHmmss"));
        map.put("service_id", "tf56pay.gateway.signQuery");
        map.put("terminal", "PC");
        map.put("version", "01");
        map.put("appid", s.getAppid());//
        map.put("sign_type", "RSA");//RSA

        map.put("businessnumber", body.getBusinessnumber());

//        String result = "";
        String result = "{\"code\":\"GP_00\",\"msg\":\"\",\"biz_code\":\"GPBIZ_00\",\"biz_msg\":\"\",\"data\":{\"sign_type\":\"RSA\",\"tf_sign\":\"1\",\"appid\":\"123\",\"businessnumber\":\"1\",\"businessrecordnumber\":\"1\",\"certcode\":\"1\",\"bankcardnumber\":\"2\",\"status\":\"1\"}}";
        try {
            String pfxFileName = s.getPfxFile();
            String keyPassword = s.getPfxSecret();
            String originalMessage = ParamUtil.sortMapByKey(map);
            System.out.println(originalMessage);
            byte[] signMessage = RSAUtils.signMessage(pfxFileName, keyPassword, originalMessage.getBytes());
            String signStr = new String(signMessage);
//            System.out.println("tf_sign:" + signStr);
            map.put("tf_sign", signStr);
            switch (sdkConfig.getMode()) {
                case "test":
                    break;
                case "normal":
                    result = HttpClient4Utils.sendHttpRequest(url, map, "utf8", true);
                    System.out.println("request:" + result);
                    break;
                default:
                    break;
            }
        } catch (Exception e) {
            e.printStackTrace();
        }

        return new SdkResp(result, "");
    }

    @PostMapping("/quickPay")
    public SdkResp quickPay(@RequestBody QuickPayReq body) {
        SdkConfig.KeyConfig s = sdkConfig.getKey(body.getAppid());
        String url = sdkConfig.getGateway();
//        System.out.printf("body=%s\n", body);

        Map<String, Object> map = new HashMap<>();
        map.put("tf_timestamp", DateUtils.strDate("yyyyMMddHHmmss"));
        map.put("service_id", "tf56pay.gateway.quickPay");
        map.put("terminal", "PC");
        map.put("version", "01");
        map.put("appid", s.getAppid());//
        map.put("sign_type", "RSA");//RSA

        map.put("backurl", body.getBackurl());
        map.put("subject", body.getSubject());
        map.put("businesstype", body.getBusinesstype());
        map.put("kind", body.getKind());
        map.put("description", body.getDescription());
        map.put("businessnumber", body.getBusinessnumber());
        map.put("billamount", body.getBillamount());
        map.put("toaccountnumber", s.getToaccountnumber());
        map.put("certcode", body.getCertcode());
        map.put("clientip", body.getClientip());
//        map.put("merchantuserid", body.getMerchantuserid());

        map.put("merchantuserid", s.getMerchantUserid());//

//        String result = "";
        String result = "{\"code\":\"GP_00\",\"msg\":\"\",\"biz_code\":\"GPBIZ_00\",\"biz_msg\":\"\",\"data\":{\"sign_type\":\"RSA\",\"tf_sign\":\"1\",\"appid\":\"123\",\"businessnumber\":\"2\",\"businessrecordnumber\":\"1\"}}";
        try {
            String pfxFileName = s.getPfxFile();
            String keyPassword = s.getPfxSecret();
            String originalMessage = ParamUtil.sortMapByKey(map);
            System.out.println(originalMessage);
            byte[] signMessage = RSAUtils.signMessage(pfxFileName, keyPassword, originalMessage.getBytes());
            String signStr = new String(signMessage);
//            System.out.println("tf_sign:" + signStr);
            map.put("tf_sign", signStr);
            switch (sdkConfig.getMode()) {
                case "test":
                    break;
                case "normal":
                    result = HttpClient4Utils.sendHttpRequest(url, map, "utf8", true);
                    System.out.println("request:" + result);
                    break;
                default:
                    break;
            }
        } catch (Exception e) {
            e.printStackTrace();
        }

        return new SdkResp(result, "");
    }

    @PostMapping("/quickPayConfirm")
    public SdkResp quickPayConfirm(@RequestBody QuickPayConfirmReq body) {
        SdkConfig.KeyConfig s = sdkConfig.getKey(body.getAppid());
        String url = sdkConfig.getGateway();
//        System.out.printf("body=%s\n", body);

        Map<String, Object> map = new HashMap<>();
        map.put("tf_timestamp", DateUtils.strDate("yyyyMMddHHmmss"));
        map.put("service_id", "tf56pay.gateway.quickPayConfirm");
        map.put("terminal", "PC");
        map.put("version", "01");
        map.put("appid", s.getAppid());//
        map.put("sign_type", "RSA");//RSA

        map.put("businessrecordnumber", body.getBusinessrecordnumber());
        map.put("verifycode", body.getVerifycode());
        map.put("clientip", body.getClientip());

//        String result = "";
        String result = "{\"code\":\"GP_00\",\"msg\":\"\",\"biz_code\":\"GPBIZ_00\",\"biz_msg\":\"\",\"data\":{\"sign_type\":\"RSA\",\"tf_sign\":\"1\",\"appid\":\"123\",\"businessnumber\":\"2\",\"businessrecordnumber\":\"1\"}}";
        try {
            String pfxFileName = s.getPfxFile();
            String keyPassword = s.getPfxSecret();
            String originalMessage = ParamUtil.sortMapByKey(map);
            System.out.println(originalMessage);
            byte[] signMessage = RSAUtils.signMessage(pfxFileName, keyPassword, originalMessage.getBytes());
            String signStr = new String(signMessage);
//            System.out.println("tf_sign:" + signStr);
            map.put("tf_sign", signStr);
            switch (sdkConfig.getMode()) {
                case "test":
                    break;
                case "normal":
                    result = HttpClient4Utils.sendHttpRequest(url, map, "utf8", true);
                    System.out.println("request:" + result);
                    break;
                default:
                    break;
            }
        } catch (Exception e) {
            e.printStackTrace();
        }

        return new SdkResp(result, "");
    }

    @PostMapping("/transfer")
    public SdkResp transfer(@RequestBody TranferReq body) {
        SdkConfig.KeyConfig s = sdkConfig.getKey(body.getAppid());
        String url = sdkConfig.getGateway();
//        System.out.printf("body=%s\n", body);

        Map<String, Object> map = new HashMap<>();
        map.put("tf_timestamp", DateUtils.strDate("yyyyMMddHHmmss"));
        map.put("service_id", "tf56enterprise.enterprise.payForCustomer");
        map.put("terminal", "PC");
        map.put("version", "01");
        map.put("appid", s.getAppid());//
        map.put("sign_type", "RSA");//RSA

        map.put("businessnumber", body.getBusinessnumber());
        map.put("subject", body.getSubject());
        map.put("transactionamount", body.getTransactionamount());
        map.put("bankcardnumber", body.getBankcardnumber());
        map.put("bankcardname", body.getBankcardname());
        map.put("bankname", body.getBankname());
        map.put("bankcardtype", "个人");//银行卡类型：个人、企业
        map.put("bankaccounttype", "储蓄卡");//银行卡借贷类型：储蓄卡、信用卡
        map.put("backurl", body.getBackurl());//付款状态通知回调地址

        map.put("fromaccountnumber", s.getToaccountnumber());//付款状态通知回调地址

//        String result = "";
        String result = "{\"code\":\"GP_00\",\"msg\":\"\",\"biz_code\":\"GPBIZ_00\",\"biz_msg\":\"\",\"data\":{\"sign_type\":\"RSA\",\"tf_sign\":\"1\",\"appid\":\"123\",\"businessnumber\":\"2\",\"businessrecordnumber\":\"1\"}}";
        try {
            String pfxFileName = s.getPfxFile();
            String keyPassword = s.getPfxSecret();
            String originalMessage = ParamUtil.sortMapByKey(map);
            System.out.println(originalMessage);
            byte[] signMessage = RSAUtils.signMessage(pfxFileName, keyPassword, originalMessage.getBytes());
            String signStr = new String(signMessage);
//            System.out.println("tf_sign:" + signStr);
            map.put("tf_sign", signStr);
            switch (sdkConfig.getMode()) {
                case "test":
                    break;
                case "normal":
                    result = HttpClient4Utils.sendHttpRequest(url, map, "utf8", true);
                    System.out.println("request:" + result);
                    break;
                default:
                    break;
            }
        } catch (Exception e) {
            e.printStackTrace();
        }

        return new SdkResp(result, "");
    }

    @PostMapping("/notify/toppay/pay")
    public Map<String, String> paycb(@RequestParam Map<String, Object> map) {
        SdkConfig.KeyConfig s = sdkConfig.getKey(map.get("appid").toString());
        String url = sdkConfig.getCbGateway();
        System.out.printf("map=%s\n", map);

        String result = "";
//        String result = "{\"code\":\"GP_00\",\"msg\":\"\",\"biz_code\":\"GPBIZ_00\",\"biz_msg\":\"\",\"data\":{\"sign_type\":\"RSA\",\"tf_sign\":\"1\",\"appid\":\"123\",\"businessnumber\":\"2\",\"businessrecordnumber\":\"1\"}}";
        try {
            String pfxFileName = s.getPfxFile();
            String keyPassword = s.getPfxSecret();
            String sign = map.get("tf_sign").toString();
            map.remove("tf_sign");
            String originalMessage = ParamUtil.sortMapByKey(map);
//            originalMessage = URLDecoder.decode(originalMessage, "UTF-8");
            System.out.println(originalMessage);
//            sign = URLDecoder.decode(sign, "UTF-8");
            System.out.println(sign);
//            byte[] signMessage = RSAUtils.signMessage(pfxFileName, keyPassword, originalMessage.getBytes());
//            String signStr = new String(signMessage);
//            System.out.println("tf_sign:" + signStr);
            boolean isOk = RSAUtils.verifySignature2(originalMessage.getBytes(StandardCharsets.UTF_8),
                    sign.getBytes(StandardCharsets.UTF_8), pfxFileName, keyPassword);
            System.out.printf("isOk=%s\n", isOk);
            map.put("sss", "123456a87g2GSG&*^Ihgqasrg");
            if (isOk) {
                result = HttpClient4Utils.sendHttpRequest(url, map, "utf8", true);
                System.out.println("request:" + result);
            }
        } catch (Exception e) {
            e.printStackTrace();
        }

        Map<String, String> ret = new TreeMap<>();
        ret.put("result", "success");
        ret.put("msg", "请求成功");
        return ret;
    }

    @PostMapping("/notify/toppay/transfer")
    public Map<String, String> transfercb(@RequestParam Map<String, Object> map) {
        SdkConfig.KeyConfig s = sdkConfig.getKey(map.get("appid").toString());
        String url = sdkConfig.getCbGateway();
        System.out.printf("map=%s\n", map);

        String result = "";
//        String result = "{\"code\":\"GP_00\",\"msg\":\"\",\"biz_code\":\"GPBIZ_00\",\"biz_msg\":\"\",\"data\":{\"sign_type\":\"RSA\",\"tf_sign\":\"1\",\"appid\":\"123\",\"businessnumber\":\"2\",\"businessrecordnumber\":\"1\"}}";
        try {
            String pfxFileName = s.getPfxFile();
            String keyPassword = s.getPfxSecret();
            String sign = map.get("tf_sign").toString();
            map.remove("tf_sign");
            String originalMessage = ParamUtil.sortMapByKey(map);
//            originalMessage = URLDecoder.decode(originalMessage, "UTF-8");
            System.out.println(originalMessage);
//            sign = URLDecoder.decode(sign, "UTF-8");
            System.out.println(sign);
//            byte[] signMessage = RSAUtils.signMessage(pfxFileName, keyPassword, originalMessage.getBytes());
//            String signStr = new String(signMessage);
//            System.out.println("tf_sign:" + signStr);
            boolean isOk = RSAUtils.verifySignature2(originalMessage.getBytes(StandardCharsets.UTF_8),
                    sign.getBytes(StandardCharsets.UTF_8), pfxFileName, keyPassword);
            System.out.printf("isOk=%s\n", isOk);
            map.put("sss", "123456a87g2GSG&*^Ihgqasrg");
            if (isOk) {
                result = HttpClient4Utils.sendHttpRequest(url, map, "utf8", true);
                System.out.println("request:" + result);
            }
        } catch (Exception e) {
            e.printStackTrace();
        }

        Map<String, String> ret = new TreeMap<>();
        ret.put("result", "success");
        ret.put("msg", "请求成功");
        return ret;
    }
}