package com.ronilsonalves.invoiceservice.data.model;

import jakarta.persistence.*;
import lombok.*;
import org.hibernate.annotations.GenericGenerator;

import java.io.Serial;
import java.io.Serializable;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.util.UUID;

@Entity
@Table(name = "invoices")
@AllArgsConstructor
@NoArgsConstructor
@Getter
@Setter
public class Invoice implements Serializable {
    @Serial
    private static final long serialVersionUID = 1912492882356572322L;

    @Id
    @GeneratedValue(generator = "system-uuid")
    @GenericGenerator(name = "system-uuid", strategy = "uuid2")
    private UUID id;

    @Temporal(TemporalType.TIMESTAMP)
    private LocalDateTime createdAt;

    private LocalDate dueDate;

    private float price;

    @Column(unique = true)
    private Integer appointmentId;

    @Temporal(TemporalType.TIMESTAMP)
    private LocalDateTime appointmentDate;

    private String appointmentDescription;

    private String patientRG;

    private String dentistCRO;
}
