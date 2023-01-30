package com.ronilsonalves.invoiceservice.data.repository;

import com.ronilsonalves.invoiceservice.data.model.Invoice;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.UUID;

@Repository
public interface InvoiceRepository extends JpaRepository<Invoice, UUID> {

    List<Invoice> getInvoicesByPatientRG(String patientRG);

    List<Invoice> getInvoicesByDentistCRO(String dentistCRO);
}