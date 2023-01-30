package com.ronilsonalves.invoiceservice.exceptionhandler;

public class UnauthorizedException extends RuntimeException {
    public UnauthorizedException(String message) {
        super(message);
    }
}
