package com.ronilsonalves.invoiceservice.data.repository;

import com.ronilsonalves.invoiceservice.config.feign.AppClientConfig;
import com.ronilsonalves.invoiceservice.data.dto.Appointment;
import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;

import java.util.List;

@FeignClient(name = "scheduling-service",
configuration = AppClientConfig.class)
public interface IAppointmentFeignRepository {

    @GetMapping("/api/v1/appointments/{appointmentId}")
    Appointment getAppointmentByID(@PathVariable Integer appointmentId);

    @GetMapping("/api/v1/appointments/patient/{patientRG}")
    List<Appointment> getAppointmentsByPatientRG(@PathVariable String patientRG);

    @GetMapping("/api/v1/appointments/dentist/{dentistCRO}")
    List<Appointment> getAppointmentsByDentistCRO(@PathVariable String dentistCRO);
}
