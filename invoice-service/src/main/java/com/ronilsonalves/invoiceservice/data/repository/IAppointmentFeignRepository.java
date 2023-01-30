package com.ronilsonalves.invoiceservice.data.repository;

import com.ronilsonalves.invoiceservice.data.dto.Appointment;
import com.ronilsonalves.invoiceservice.data.dto.Patient;
import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;

import java.util.List;

@FeignClient(name = "schedule-service")
public interface IAppointmentFeignRepository {

    @GetMapping("/appointments/{appointmentId}")
    Appointment getAppointmentByID(@PathVariable Integer appointmentId);

    @GetMapping("/appointments/patient/{patientRG}")
    List<Appointment> getAppointmentsByPatientRG(@PathVariable String patientRG);

    @GetMapping("/appointments/dentist/{dentistCRO}")
    List<Appointment> getAppointmentsByDentistCRO(@PathVariable String dentistCRO);
}
