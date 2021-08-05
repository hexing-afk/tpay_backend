package com.example.model;

public class TranferReq {
    private String appid;
    private String businessnumber;//业务流水号	商户端的流水号，需保证在商户端不重复
    private String subject;//商品名称	商品名称
    private String transactionamount;//交易金额	单位：元
    private String bankcardnumber;//银行卡号码
    private String bankcardname;//银行卡姓名
    private String bankname;//银行名称	银行名称参考附录1里标准银行名称
    private String backurl;//回调客户端url	付款状态通知回调地址

    public TranferReq(String appid) {
        this.appid = appid;
    }

    public String getBusinessnumber() {
        return businessnumber;
    }

    public void setBusinessnumber(String businessnumber) {
        this.businessnumber = businessnumber;
    }

    public String getSubject() {
        return subject;
    }

    public void setSubject(String subject) {
        this.subject = subject;
    }

    public String getTransactionamount() {
        return transactionamount;
    }

    public void setTransactionamount(String transactionamount) {
        this.transactionamount = transactionamount;
    }

    public String getBankcardnumber() {
        return bankcardnumber;
    }

    public void setBankcardnumber(String bankcardnumber) {
        this.bankcardnumber = bankcardnumber;
    }

    public String getBankcardname() {
        return bankcardname;
    }

    public void setBankcardname(String bankcardname) {
        this.bankcardname = bankcardname;
    }

    public String getBankname() {
        return bankname;
    }

    public void setBankname(String bankname) {
        this.bankname = bankname;
    }

    public String getBackurl() {
        return backurl;
    }

    public void setBackurl(String backurl) {
        this.backurl = backurl;
    }

    public String getAppid() {
        return appid;
    }

    public void setAppid(String appid) {
        this.appid = appid;
    }
}
