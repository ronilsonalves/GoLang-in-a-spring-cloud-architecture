package com.ronilsonalves.invoiceservice.data.dto;

import java.time.LocalDateTime;

public record Patient(Integer id, String name, String lastName, String rg, LocalDateTime createdAt) {
}
