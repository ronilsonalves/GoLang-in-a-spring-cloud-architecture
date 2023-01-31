package com.ronilsonalves.invoiceservice.api.service.impl;

import com.ronilsonalves.invoiceservice.api.service.InvoiceService;
import com.ronilsonalves.invoiceservice.data.dto.Appointment;
import com.ronilsonalves.invoiceservice.data.dto.InvoiceRequestBody;
import com.ronilsonalves.invoiceservice.data.model.Invoice;
import com.ronilsonalves.invoiceservice.data.repository.IAppointmentFeignRepository;
import com.ronilsonalves.invoiceservice.data.repository.InvoiceRepository;
import com.ronilsonalves.invoiceservice.exceptionhandler.ResourceNotFoundException;
import lombok.RequiredArgsConstructor;
import org.springframework.beans.BeanUtils;
import org.springframework.stereotype.Service;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Service
@RequiredArgsConstructor
public class InvoiceServiceImpl implements InvoiceService {
    private final InvoiceRepository repository;
    private final IAppointmentFeignRepository appointmentFeignRepository;

    /**
     * Returns a list of invoices, if there no invoices the list returned will be empty
     * @return a list of invoices, if there no invoices the list returned will be empty
     */
    @Override
    public List<Invoice> listAll() {
        return repository.findAll();
    }

    /**
     * Returns an Invoice with ID generated by database application
     * @param invoice an invoiceRequestBody object provided to save at database
     * @return an Invoice with ID generated by database application
     * @throws ResourceNotFoundException if there are no appointments with
     * ID provided.
     */
    @Override
    public Invoice create(InvoiceRequestBody invoice) {
        Invoice invoiceToSave = new Invoice();
        BeanUtils.copyProperties(invoice,invoiceToSave);
        Appointment appointment = appointmentFeignRepository.getAppointmentByID(invoice.getAppointmentId());
        if (appointment != null) {
            return repository.save(buildInvoiceToSave(invoiceToSave,appointment));
        } else throw new ResourceNotFoundException("Appointment's ID provided was not found. " +
                "Please check the ID provided");
    }

    /**
     * Returns the object with data updated
     * @param invoiceRequestBody an invoice DTO object to update data of invoice with UUID provided
     * @param uuid invoice identification number
     * @return the updated invoice with UUID provided
     * @throws ResourceNotFoundException if the UUID provided was not found at database
     * @throws com.ronilsonalves.invoiceservice.exceptionhandler.InvalidUUIDException if an invalid UUID is given
     */
    @Override
    public Invoice update(InvoiceRequestBody invoiceRequestBody, UUID uuid) {
        Invoice toUpdate = new Invoice();
        Appointment appointment = appointmentFeignRepository.getAppointmentByID(invoiceRequestBody.getAppointmentId());
        if (repository.findById(uuid).isPresent()) {
            BeanUtils.copyProperties(invoiceRequestBody,toUpdate);
        } else throw new ResourceNotFoundException("There's no invoice with UUID "+uuid+" provided");
        if (appointment != null) {
            return repository.save(buildInvoiceToSave(toUpdate,appointment));
        } else throw new ResourceNotFoundException("Appointment's ID provided was not found. " +
                "Please check the ID provided");
    }

    /**
     * Returns an invoice with provided ID
     * @param invoiceId key to find the invoice at database
     * @return an invoice with provided ID
     * @throws ResourceNotFoundException if invoiceId provided is
     * not found at database
     */
    @Override
    public Invoice getInvoiceByID(UUID invoiceId) {
        Optional<Invoice> response = repository.findById(invoiceId);
        if (response.isPresent()) {
            return response.get();
        } else throw new ResourceNotFoundException("There no invoices with ID "+invoiceId+" provided.");
    }

    /**
     * Returns a list of invoices from a patient from their RG document
     * @param patientRG brazilian ID document used to verify if the patient have any invoice registered at database
     * @return a list of invoices from a patient from their RG document
     * @throws ResourceNotFoundException if patientRG provided not
     * found at database
     */
    @Override
    public List<Invoice> getInvoicesByPatientRG(String patientRG) {
        List<Invoice> response = repository.getInvoicesByPatientRG(patientRG);
        if (response != null && response.size() > 0) {
            return response;
        } else throw new ResourceNotFoundException("There no invoices for patient with RG "+patientRG+" provided");
    }

    /**
     * Returns a list of invoices where a dentist was the provider of dental services
     * @param dentistCRO brazilian profession council ID document used to verify if the dentist have any invoice
     *                   registered at database as provider
     * @return a list of invoices where a dentist was the provider of dental services
     * @throws ResourceNotFoundException if dentistCRO provided not found at database
     */
    @Override
    public List<Invoice> getInvoicesByDentistCRO(String dentistCRO) {
        List<Invoice> response = repository.getInvoicesByDentistCRO(dentistCRO);
        if (response != null && !response.isEmpty()) {
            return response;
        } else throw new ResourceNotFoundException("There no invoices where dentist with CRO "+dentistCRO+
                " is provider of dental services");
    }

    /**
     * Create a new invoice from RabbitMQ queue
     * @param appointment the appointment record provided by the RabbitMQ queue
     */
    public void save(Appointment appointment) {
        Invoice invoiceToSave = new Invoice();
        repository.save(buildInvoiceToSave(invoiceToSave,appointment));
    }

    /**
     * @param invoiceId used to verify if exists an invoice registered at database
     * @throws ResourceNotFoundException if invoiceID provided there no exist at database
     */
    @Override
    public void delete(UUID invoiceId) {
        if (repository.findById(invoiceId).isPresent()) {
            repository.deleteById(invoiceId);
        } else throw new ResourceNotFoundException("Unable to delete. An invoice with ID "+invoiceId+" was not found");
    }

    /**
     * Return an invoice object built from invoiceRequestBody and appointment data
     * @param invoice the parsed object from invoiceRequestBody
     * @param appointment the appointment record from method's call
     * @return an invoice object built from invoiceRequestBody and appointment data
     */
    private Invoice buildInvoiceToSave(Invoice invoice,Appointment appointment) {
        invoice.setCreatedAt(LocalDateTime.now());
        invoice.setDueDate(LocalDate.from(appointment.dateAndTime().plusDays(5)));
        invoice.setAppointmentId(appointment.id());
        invoice.setAppointmentDate(appointment.dateAndTime());
        invoice.setAppointmentDescription(appointment.description());
        invoice.setPatientRG(appointment.patient().rg());
        invoice.setDentistCRO(appointment.dentist().cro());
        // setting default price to R$ 99,00 for the appointment data come from micro service via RabbitMQ
        if (invoice.getPrice() == 0) {
            invoice.setPrice(99);
        }
        return invoice;
    }
}