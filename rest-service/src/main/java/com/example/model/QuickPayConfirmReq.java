package com.example.model;

public class QuickPayConfirmReq {
    private String appid;
    private String businessrecordnumber;
    private String verifycode;
    private String clientip;

    public String getBusinessrecordnumber() {
        return businessrecordnumber;
    }

    public void setBusinessrecordnumber(String businessrecordnumber) {
        this.businessrecordnumber = businessrecordnumber;
    }

    public String getVerifycode() {
        return verifycode;
    }

    public void setVerifycode(String verifycode) {
        this.verifycode = verifycode;
    }

    public String getClientip() {
        return clientip;
    }

    public void setClientip(String clientip) {
        this.clientip = clientip;
    }

    public String getAppid() {
        return appid;
    }

    public void setAppid(String appid) {
        this.appid = appid;
    }
}
