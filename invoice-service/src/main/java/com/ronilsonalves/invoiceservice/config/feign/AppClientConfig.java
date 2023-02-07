package com.ronilsonalves.invoiceservice.config.feign;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.ronilsonalves.invoiceservice.config.CustomLocalDateTimeObjectMapper;
import feign.codec.Decoder;
import feign.jackson.JacksonDecoder;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class AppClientConfig {
    @Bean
    public Decoder feignDecoder() {
        CustomLocalDateTimeObjectMapper customLocalDateTimeObjectMapper =
                new CustomLocalDateTimeObjectMapper(new ObjectMapper());
        return new JacksonDecoder(customLocalDateTimeObjectMapper.getCustomLocalDateTimeObjectMapper());
    }
}
