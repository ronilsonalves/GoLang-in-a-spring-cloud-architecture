package com.ronilsonalves.invoiceservice.queue;

import com.ronilsonalves.invoiceservice.api.service.impl.InvoiceServiceImpl;
import com.ronilsonalves.invoiceservice.data.dto.Appointment;
import lombok.RequiredArgsConstructor;
import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.stereotype.Component;

import java.util.logging.Level;
import java.util.logging.Logger;

@RequiredArgsConstructor
@Component
public class AppointmentListener {
    private final Logger logger = Logger.getLogger(AppointmentListener.class.getName());
    private final InvoiceServiceImpl invoiceService;

    @RabbitListener(queues = {"${queue.appointment-service.name}"})
    public void receiveMessage(Appointment appointment) {
        logger.log(Level.INFO,"Message received from RabbitMQ: "+appointment);
        invoiceService.save(appointment);
    }
}
