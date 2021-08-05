package com.example.restservice;

public class SdkResp {
    private final String resp;
    private String error;

    public SdkResp(String resp, String error) {
        this.resp = resp;
        this.error = error;
    }

    public String getResp() {
        return resp;
    }

    public String getError() {
        return error;
    }

    public void setError(String error) {
        this.error = error;
    }
}