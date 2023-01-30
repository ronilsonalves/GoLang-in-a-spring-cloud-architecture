package com.ronilsonalves.invoiceservice.data.repository;

import com.ronilsonalves.invoiceservice.data.dto.Patient;
import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;

@FeignClient(name = "schedule-service")
public interface IPatientFeignRepository {

    @GetMapping("/patients/RG/{patientRG}")
    Patient getPatientByRG(@PathVariable String patientRG);
}
