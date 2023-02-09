package com.ronilsonalves.invoiceservice.api.controller;

import com.ronilsonalves.invoiceservice.api.service.InvoiceService;
import com.ronilsonalves.invoiceservice.data.dto.InvoiceRequestBody;
import com.ronilsonalves.invoiceservice.data.model.Invoice;
import com.ronilsonalves.invoiceservice.exceptionhandler.InternalServerException;
import com.ronilsonalves.invoiceservice.exceptionhandler.InvalidUUIDException;
import com.ronilsonalves.invoiceservice.exceptionhandler.ResourceNotFoundException;
import com.ronilsonalves.invoiceservice.exceptionhandler.UnauthorizedException;
import io.swagger.v3.oas.annotations.media.ArraySchema;
import io.swagger.v3.oas.annotations.media.Content;
import io.swagger.v3.oas.annotations.media.Schema;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.UUID;

@RestController
@RequestMapping("/invoices")
@RequiredArgsConstructor
@Tag(name = "Invoices")
public class InvoiceController {

    private final InvoiceService service;


    @ApiResponses(value = {
            @ApiResponse(
                    responseCode = "200",
                    description = "Returns a list unengaged of invoices",
                    content = {
                            @Content(
                                    mediaType = "application/json",
                                    array = @ArraySchema(schema = @Schema(implementation = Invoice.class))
                            )
                    })
    })
    @SecurityRequirement(name = "oauth_sec")
    @GetMapping
    public ResponseEntity<?> getAllInvoices() throws UnauthorizedException{
        return ResponseEntity.ok().body(service.listAll());
    }

    @ApiResponses(value = {
            @ApiResponse(
                    responseCode = "201",
                    description = "Returns the created invoice",
                    content = {
                            @Content(
                                    mediaType = "application/json",
                                    schema = @Schema(implementation = Invoice.class)
                            )
                    })
    })
    @SecurityRequirement(name = "oauth_sec")
    @PostMapping
    public ResponseEntity<Invoice> save(@RequestBody InvoiceRequestBody invoiceRequestBody) throws InvalidUUIDException,
            InternalServerException, UnauthorizedException {
        return ResponseEntity.status(HttpStatus.CREATED).body(service.create(invoiceRequestBody));
    }

    @ApiResponses(value = {
            @ApiResponse(
                    responseCode = "200",
                    description = "Returns the updated invoice",
                    content = {
                            @Content(
                                    mediaType = "application/json",
                                    schema = @Schema(implementation = Invoice.class)
                            )
                    })
    })
    @SecurityRequirement(name = "oauth_sec")
    @PutMapping("/{invoiceID}")
    public ResponseEntity<Invoice> update(@PathVariable String invoiceID,
                                          @RequestBody InvoiceRequestBody invoiceRequestBody) throws InvalidUUIDException,
            InternalServerException, UnauthorizedException{
        return ResponseEntity.status(HttpStatus.OK).body(service.update(invoiceRequestBody, UUID.fromString(invoiceID)));
    }

    @ApiResponses(value = {
            @ApiResponse(
                    responseCode = "200",
                    description = "Return an invoice by Id",
                    content = {
                            @Content(
                                    mediaType = "application/json",
                                    schema = @Schema(implementation = Invoice.class)
                            )
                    }),
            @ApiResponse(
                    responseCode = "400",
                    description = "Throws a bad request exception. Caused by an invalid UUID provided at request",
                    content = {
                            @Content(
                                    mediaType = "application/json"
                            )
                    }),
            @ApiResponse(
                    responseCode = "404",
                    description = "Throws a NOT FOUND exception. An invoice with provided ID was not found",
                    content = {
                            @Content(
                                    mediaType = "application/json"
                            )
                    })
    })
    @SecurityRequirement(name = "oauth_sec")
    @GetMapping("/{invoiceID}")
    public ResponseEntity<Invoice> getInvoiceById(@PathVariable String invoiceID) throws InvalidUUIDException,
            ResourceNotFoundException,
            UnauthorizedException {
        try {
            return ResponseEntity.ok().body(service.getInvoiceByID(UUID.fromString(invoiceID)));
        } catch (IllegalArgumentException e) {
            throw new InvalidUUIDException("The UUID provided is not in a valid format, please check the UUID and try" +
                    " again.");
        }

    }

    @ApiResponses(value = {
            @ApiResponse(
                    responseCode = "200",
                    description = "Returns a list of invoices from a patient from their RG document",
                    content = {
                            @Content(
                                    mediaType = "application/json",
                                    schema = @Schema(implementation = Invoice.class)
                            )
                    }),
            @ApiResponse(
                    responseCode = "400",
                    description = "Throws a InvalidUUID exception. The UUID provided is invalid ",
                    content = {
                            @Content(mediaType = "application/json")
                    }
            ),
            @ApiResponse(
                    responseCode = "404",
                    description = "Throws a NOT FOUND exception. There no invoices registered for the patient with " +
                            "the provided RG document or there no patient with the provided RG",
                    content = {
                            @Content(mediaType = "application/json")
                    }
            )
    })
    @SecurityRequirement(name = "oauth_sec")
    @GetMapping("/patient/{patientRG}")
    public ResponseEntity<List<Invoice>> getInvoicesByPatientRG(@PathVariable String patientRG) throws ResourceNotFoundException, UnauthorizedException {
        return ResponseEntity.ok().body(service.getInvoicesByPatientRG(patientRG));
    }

    @ApiResponses(value = {
            @ApiResponse(
                    responseCode = "200",
                    description = "Returns a list of invoices where a dentist was the provider of dental services",
                    content = {
                            @Content(
                                    mediaType = "application/json",
                                    array = @ArraySchema(schema = @Schema(implementation = Invoice.class))
                            )
                    }),
            @ApiResponse(
                    responseCode = "404",
                    description = "Throws a NOT FOUND exception. There no invoices registered with the dentistCRO " +
                            "informed or there no dentist with the provided CRO document",
                    content = {
                            @Content(mediaType = "application/json")
                    }
            )
    })
    @SecurityRequirement(name = "oauth_sec")
    @GetMapping("/dentist/{dentistCRO}")
    public ResponseEntity<List<Invoice>> getInvoicesByDentistCRO(@PathVariable String dentistCRO) throws ResourceNotFoundException, UnauthorizedException {
        return ResponseEntity.ok().body(service.getInvoicesByDentistCRO(dentistCRO));
    }

    @ApiResponses(value = {
            @ApiResponse(
                    responseCode = "200",
                    description = "Delete an invoice successfully",
                    content = {
                            @Content(
                                    mediaType = "application/json",
                                    array = @ArraySchema(schema = @Schema(implementation = Invoice.class))
                            )
                    }),
            @ApiResponse(
                    responseCode = "404",
                    description = "Throws a NOT FOUND exception. There no invoices registered with the invoiceID " +
                            "informed",
                    content = {
                            @Content(mediaType = "application/json")
                    }
            )
    })
    @SecurityRequirement(name = "oauth_sec")
    @DeleteMapping("/{invoiceID}")
    public ResponseEntity<?> deleteInvoiceByID(@PathVariable String invoiceID) throws ResourceNotFoundException,
            UnauthorizedException {
        service.delete(UUID.fromString(invoiceID));
        return ResponseEntity.ok().body("Invoice deleted successfully");
    }
}
