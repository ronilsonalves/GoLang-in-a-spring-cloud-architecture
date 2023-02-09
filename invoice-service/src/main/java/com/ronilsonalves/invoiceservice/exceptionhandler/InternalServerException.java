package com.ronilsonalves.invoiceservice.exceptionhandler;

public class InternalServerException extends RuntimeException {
    public InternalServerException(String message) {
        super(message);
    }
}
