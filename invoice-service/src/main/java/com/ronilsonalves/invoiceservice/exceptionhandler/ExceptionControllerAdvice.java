package com.ronilsonalves.invoiceservice.exceptionhandler;

import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestControllerAdvice;
import org.springframework.web.servlet.mvc.method.annotation.ResponseEntityExceptionHandler;

import java.time.LocalDateTime;

@RestControllerAdvice(basePackages = "com.ronilsonalves.invoiceservice.api.controller")
public class ExceptionControllerAdvice extends ResponseEntityExceptionHandler {

    @ResponseBody
    @ExceptionHandler(ResourceNotFoundException.class)
    @ResponseStatus(HttpStatus.NOT_FOUND)
    public ResponseEntity<MessageExceptionHandler> resourceNotFoundException(ResourceNotFoundException resourceNotFoundException) {
        MessageExceptionHandler error = new MessageExceptionHandler(
                LocalDateTime.now(), HttpStatus.NOT_FOUND.value(), resourceNotFoundException.getMessage()
        );
        return new ResponseEntity<>(error,HttpStatus.NOT_FOUND);
    }

    @ResponseBody
    @ExceptionHandler(UnauthorizedException.class)
    @ResponseStatus(HttpStatus.UNAUTHORIZED)
    public ResponseEntity<MessageExceptionHandler> unauthorizedException(UnauthorizedException unauthorizedException) {
        MessageExceptionHandler error = new MessageExceptionHandler(
                LocalDateTime.now(), HttpStatus.UNAUTHORIZED.value(), unauthorizedException.getMessage()
        );
        return new ResponseEntity<>(error,HttpStatus.UNAUTHORIZED);
    }

    @ResponseBody
    @ExceptionHandler(InvalidUUIDException.class)
    @ResponseStatus(HttpStatus.BAD_REQUEST)
    public ResponseEntity<MessageExceptionHandler> invalidUUIDException(InvalidUUIDException invalidUUIDException) {
        MessageExceptionHandler error = new MessageExceptionHandler(
                LocalDateTime.now(), HttpStatus.BAD_REQUEST.value(), invalidUUIDException.getMessage()
        );
        return new ResponseEntity<>(error,HttpStatus.BAD_REQUEST);
    }

}