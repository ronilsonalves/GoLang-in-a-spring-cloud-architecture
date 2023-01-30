package com.ronilsonalves.invoiceservice.exceptionhandler;

public class InvalidUUIDException extends IllegalArgumentException{
    public InvalidUUIDException(String message) {
        super(message);
    }
}
