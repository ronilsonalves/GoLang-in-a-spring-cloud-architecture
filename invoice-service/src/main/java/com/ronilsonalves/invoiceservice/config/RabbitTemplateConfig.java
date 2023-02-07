package com.ronilsonalves.invoiceservice.config;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.amqp.rabbit.connection.ConnectionFactory;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.amqp.support.converter.Jackson2JsonMessageConverter;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class RabbitTemplateConfig {

    @Bean
    public Jackson2JsonMessageConverter producerJackson2MessageConverter() {
        CustomLocalDateTimeObjectMapper customLocalDateTimeObjectMapper =
                new CustomLocalDateTimeObjectMapper(new ObjectMapper());
//        ObjectMapper om = new ObjectMapper();
//        JavaTimeModule module = new JavaTimeModule();
//        om.setVisibility(PropertyAccessor.ALL, JsonAutoDetect.Visibility.ANY);
//        om.disable(SerializationFeature.WRITE_DATES_AS_TIMESTAMPS);
//        LocalDateTimeDeserializer localDateTimeDeserializer =
//                new LocalDateTimeDeserializer(DateTimeFormatter.ofPattern("dd/MM/yyyy HH:mm"));
//        module.addDeserializer(LocalDateTime.class, localDateTimeDeserializer);
//        om.registerModule(module);
        return new Jackson2JsonMessageConverter(customLocalDateTimeObjectMapper.getCustomLocalDateTimeObjectMapper());
    }

    @Bean
    public RabbitTemplate rabbitTemplate(ConnectionFactory connectionFactory) {
        RabbitTemplate rabbitTemplate = new RabbitTemplate(connectionFactory);
        rabbitTemplate.setMessageConverter(producerJackson2MessageConverter());
        return rabbitTemplate;
    }
}
