package com.ronilsonalves.invoiceservice.config;

import com.fasterxml.jackson.annotation.JsonAutoDetect;
import com.fasterxml.jackson.annotation.PropertyAccessor;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.fasterxml.jackson.datatype.jsr310.deser.LocalDateTimeDeserializer;
import lombok.RequiredArgsConstructor;

import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;

@RequiredArgsConstructor
public class CustomLocalDateTimeObjectMapper {
    private final ObjectMapper customObjectMapper;

    public ObjectMapper getCustomLocalDateTimeObjectMapper() {
        JavaTimeModule module = new JavaTimeModule();
        this.customObjectMapper.setVisibility(PropertyAccessor.ALL, JsonAutoDetect.Visibility.ANY);
        this.customObjectMapper.disable(SerializationFeature.WRITE_DATES_AS_TIMESTAMPS);
        LocalDateTimeDeserializer localDateTimeDeserializer =
                new LocalDateTimeDeserializer(DateTimeFormatter.ofPattern("dd/MM/yyyy HH:mm"));
        module.addDeserializer(LocalDateTime.class, localDateTimeDeserializer);
        this.customObjectMapper.registerModule(module);
        return this.customObjectMapper;
    }
}
