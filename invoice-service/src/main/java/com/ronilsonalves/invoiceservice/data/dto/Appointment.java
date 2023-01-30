package com.ronilsonalves.invoiceservice.data.dto;

import java.time.LocalDateTime;

public record Appointment (Integer id, String description, LocalDateTime dateAndTime, String dentistCRO,
                           String patientRG, Dentist dentist, Patient patient) {

}