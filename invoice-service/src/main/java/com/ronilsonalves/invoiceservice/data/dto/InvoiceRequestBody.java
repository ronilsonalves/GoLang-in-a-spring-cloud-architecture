package com.ronilsonalves.invoiceservice.data.dto;

import jakarta.validation.constraints.NotEmpty;
import lombok.Data;
import lombok.Getter;
import lombok.Setter;

@Data
@Getter
@Setter
public class InvoiceRequestBody {

    @NotEmpty
    private float price;

    @NotEmpty
    private Integer appointmentId;
}
