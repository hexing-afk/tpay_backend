package com.example.model;

public class QuickPayReq {
    private String appid;
    private String backurl;
    private String subject;
    private String businesstype;
    private String kind;
    private String description;
    private String businessnumber;
    private String billamount;
    private String toaccountnumber;
    private String certcode;
    private String clientip;
    private String merchantuserid;

    public String getBackurl() {
        return backurl;
    }

    public void setBackurl(String backurl) {
        this.backurl = backurl;
    }

    public String getSubject() {
        return subject;
    }

    public void setSubject(String subject) {
        this.subject = subject;
    }

    public String getBusinesstype() {
        return businesstype;
    }

    public void setBusinesstype(String businesstype) {
        this.businesstype = businesstype;
    }

    public String getKind() {
        return kind;
    }

    public void setKind(String kind) {
        this.kind = kind;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public String getBusinessnumber() {
        return businessnumber;
    }

    public void setBusinessnumber(String businessnumber) {
        this.businessnumber = businessnumber;
    }

    public String getBillamount() {
        return billamount;
    }

    public void setBillamount(String billamount) {
        this.billamount = billamount;
    }

    public String getToaccountnumber() {
        return toaccountnumber;
    }

    public void setToaccountnumber(String toaccountnumber) {
        this.toaccountnumber = toaccountnumber;
    }

    public String getCertcode() {
        return certcode;
    }

    public void setCertcode(String certcode) {
        this.certcode = certcode;
    }

    public String getClientip() {
        return clientip;
    }

    public void setClientip(String clientip) {
        this.clientip = clientip;
    }

    public String getMerchantuserid() {
        return merchantuserid;
    }

    public void setMerchantuserid(String merchantuserid) {
        this.merchantuserid = merchantuserid;
    }

    public String getAppid() {
        return appid;
    }

    public void setAppid(String appid) {
        this.appid = appid;
    }
}
