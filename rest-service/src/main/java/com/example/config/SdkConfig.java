package com.example.config;

import lombok.Data;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.PropertySource;

import java.util.ArrayList;
import java.util.List;
import java.util.Optional;

@Data
@Configuration
@ConfigurationProperties(prefix = "spring.sdk")
@PropertySource(value = "application.yml", encoding = "utf-8")
public class SdkConfig {
    private String gateway;
    private String mode;
    private String cbGateway;
    private List<KeyConfig> keys = new ArrayList<>();

    public KeyConfig getKey(String key) {
        Optional<KeyConfig> f = keys.stream().filter(p -> key.equals(p.appid)).findFirst();
        return f.orElse(null);
    }

    @Data
    public static class KeyConfig {
        private String pfxFile;
        private String pfxSecret;
        private String appid;
        private String merchantUserid;
        private String toaccountnumber;
    }
}
