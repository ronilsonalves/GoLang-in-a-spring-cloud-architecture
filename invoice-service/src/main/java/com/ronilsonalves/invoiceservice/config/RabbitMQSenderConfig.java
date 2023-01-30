package com.ronilsonalves.invoiceservice.config;

import org.springframework.amqp.core.Queue;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class RabbitMQSenderConfig {
    @Value("${queue.appointment-service.name}")
    private String appointmentServiceQueue;

    @Bean
    public Queue appointmentQueue() {
        return new Queue(this.appointmentServiceQueue,false);
    }
}
