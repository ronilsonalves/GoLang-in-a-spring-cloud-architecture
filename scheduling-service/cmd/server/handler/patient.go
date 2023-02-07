package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ronilsonalves/GoLang-in-a-spring-cloud-architecture/scheduling-service/internal/domain"
	"github.com/ronilsonalves/GoLang-in-a-spring-cloud-architecture/scheduling-service/internal/patient"
	"github.com/ronilsonalves/GoLang-in-a-spring-cloud-architecture/scheduling-service/pkg/web"
	"net/http"
	"strconv"
)

type patientHandler struct {
	s patient.Service
}

func NewPatientHandler(s patient.Service) *patientHandler {
	return &patientHandler{
		s: s,
	}
}

// GetAll - get all patients from db.
// @BasePath /api/v1
// GetAllPatients godoc
// @Summary List all patients
// @Schemes
// @Description get all patients from db.
// @Tags Patients
// @Accept json
// @Produce json
// @Success 200 {object} []domain.Patient
// @Failure 400 {object} web.errorResponse
// @Failure 401 {object} web.errorResponse
// @Router /patients [get]
// @Security OAuth2Application
func (h *patientHandler) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		patients, err := h.s.GetAll()
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, patients)
	}
}

// GetByID - get a patient by an ID
// @BasePath /api/v1
// GetPatientByID godoc
// @Summary Get a patient by an ID
// @Schemes
// @Description get a patient by a provided ID.
// @Tags Patients
// @Accept json
// @Produce json
// @Param id path int true "Patient ID"
// @Success 200 {object} domain.Patient
// @Failure 400 {object} web.errorResponse
// @Failure 401 {object} web.errorResponse
// @Failure 404 {object} web.errorResponse
// @Router /patients/{id} [get]
// @Security OAuth2Application
func (h *patientHandler) GetByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid id provided")
			return
		}

		patient, err := h.s.GetByID(id)
		if err != nil {
			web.BadResponse(ctx, http.StatusNotFound, "error", "patient not found")
			return
		}
		web.ResponseOK(ctx, http.StatusOK, patient)
	}
}

// Post - save a new patient into db
// @BasePath /api/v1
// PostPatient godoc
// @Summary Create a new patient
// @Schemes
// @Description Create a new patient by request body.
// @Tags Patients
// @Accept json
// @Produce json
// @Param body body domain.Patient true "Body"
// @Success 201 {object} domain.Patient
// @Failure 400 {object} web.errorResponse
// @Failure 401 {object} web.errorResponse
// @Router /patients [post]
// @Security OAuth2Application
func (h *patientHandler) Post() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var patient domain.Patient
		err := ctx.ShouldBindJSON(&patient)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid patient")
			return
		}

		isValid, err := isEmptyPatient(&patient)
		if !isValid {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}

		response, err := h.s.Create(patient)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}

		web.ResponseOK(ctx, http.StatusCreated, response)
	}
}

// Put - update an entire patient
// @BasePath /api/v1
// PutPatient godoc
// @Summary Update an entire patient by ID
// @Schemes
// @Description Update an entire patient by ID.
// @Tags Patients
// @Accept json
// @Produce json
// @Param id path int true "Patient ID"
// @Param body body domain.Patient true "Body"
// @Success 200 {object} domain.Patient
// @Failure 400 {object} web.errorResponse
// @Failure 401 {object} web.errorResponse
// @Failure 409 {object} web.errorResponse
// @Router /patients [put]
// @Security OAuth2Application
func (h *patientHandler) Put() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid patient id provided")
			return
		}
		var patient domain.Patient
		err = ctx.ShouldBindJSON(&patient)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid patient data")
		}

		isValid, err := isEmptyPatient(&patient)
		if !isValid {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}

		response, err := h.s.Update(id, patient)
		if err != nil {
			web.BadResponse(ctx, http.StatusConflict, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}

// Patch - update fields from a patient
// @BasePath /api/v1
// PatchPatient godoc
// @Summary Update fields from a patient
// @Schemes
// @Description Update fields from a patient
// @Tags Patients
// @Accept json
// @Produce json
// @Param id path int true "Patient ID"
// @Param body body domain.Patient true "Body"
// @Success 200 {object} domain.Patient
// @Failure 400 {object} web.errorResponse
// @Failure 401 {object} web.errorResponse
// @Router /patients/{id} [patch]
// @Security OAuth2Application
func (h *patientHandler) Patch() gin.HandlerFunc {
	type Request struct {
		Surname        string `json:"surname,omitempty"`
		Name           string `json:"name,omitempty"`
		IdentityNumber string `json:"identity_number,omitempty"`
		CreatedAt      string `json:"created_at,omitempty"`
	}
	return func(ctx *gin.Context) {
		var r Request
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid id provided")
			return
		}
		if err := ctx.ShouldBindJSON(&r); err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid request")
			return
		}
		update := domain.Patient{
			LastName:  r.Surname,
			Name:      r.Name,
			RG:        r.IdentityNumber,
			CreatedAt: r.CreatedAt,
		}
		if update.CreatedAt != "" {
			if !validateDateTime(update.CreatedAt) {
				web.BadResponse(ctx, http.StatusBadRequest, "error", "please the patient created_at field must be in format: 30/01/2023 23:59 or 30/01/2023 23:59:59")
				return
			}
		}
		response, err := h.s.Update(id, update)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}

// Delete - delete a patient
// @BasePath /api/v1
// DeletePatient godoc
// @Summary Delete a patient by ID
// @Schemes
// @Description Delete a patient by ID
// @Tags Patients
// @Accept json
// @Produce json
// @Param id path int true "Patient ID"
// @Success 200 {object} web.response
// @Failure 400 {object} web.errorResponse
// @Failure 404 {object} web.errorResponse
// @Failure 401 {object} web.errorResponse
// @Router /patients/{id} [delete]
// @Security OAuth2Application
func (h *patientHandler) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid id provided")
			return
		}
		err = h.s.Delete(id)
		if err != nil {
			web.BadResponse(ctx, http.StatusNotFound, "error", err.Error())
			return
		}
		web.DeleteResponse(ctx, http.StatusOK, "patient deleted")
	}
}

// Aux functions bellow ->

func isEmptyPatient(patient *domain.Patient) (bool, error) {
	switch {
	case patient.LastName == "" || patient.Name == "" || patient.CreatedAt == "" || patient.RG == "":
		return false, errors.New("patient fields can't be empty")
	case !validateDateTime(patient.CreatedAt):
		return false, errors.New("please the patient created_at field must be in format: 30/01/2023 23:59 or 30/01/2023 23:59:59")
	}
	return true, nil
}
