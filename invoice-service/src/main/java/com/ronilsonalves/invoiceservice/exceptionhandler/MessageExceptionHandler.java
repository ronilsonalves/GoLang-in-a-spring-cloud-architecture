package com.ronilsonalves.invoiceservice.exceptionhandler;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.Setter;

import java.time.LocalDateTime;

@AllArgsConstructor
@Getter
@Setter
public class MessageExceptionHandler {
    private LocalDateTime timeStamp;
    private Integer status;
    private String message;
}
