package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ronilsonalves/GoLang-in-a-spring-cloud-architecture/scheduling-service/internal/appointment"
	"github.com/ronilsonalves/GoLang-in-a-spring-cloud-architecture/scheduling-service/internal/domain"
	"github.com/ronilsonalves/GoLang-in-a-spring-cloud-architecture/scheduling-service/pkg/web"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type appointmentHandler struct {
	s appointment.Service
}

func NewAppointmentHandler(s appointment.Service) *appointmentHandler {
	return &appointmentHandler{
		s: s,
	}
}

// GetAll - get all appointments from db.
// @BasePath /api/v1
// GetAllAppointments godoc
// @Summary List all appointments
// @Schemes
// @Description get all appointments from db.
// @Tags Appointments
// @Accept json
// @Produce json
// @Success 200 {object} []domain.AppointmentDTO
// @Failure 400 {object} web.errorResponse
// @Failure 401 {object} web.errorResponse
// @Failure 404 {object} web.errorResponse
// @Router /appointments [get]
// @Security OAuth2Application
func (h *appointmentHandler) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		response, err := h.s.GetAll()
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}
		if response == nil {
			web.BadResponse(ctx, http.StatusNotFound, "error", "was not found appointments registered")
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}

// GetByID - get a appointment by an ID
// @BasePath /api/v1
// GetAppointmentByID godoc
// @Summary Get an appointment by an ID
// @Schemes
// @Description get an appointment by a provided ID.
// @Tags Appointments
// @Accept json
// @Produce json
// @Param id path int true "Appointment ID"
// @Success 200 {object} domain.AppointmentDTO
// @Failure 400 {object} web.errorResponse
// @Failure 401 {object} web.errorResponse
// @Failure 404 {object} web.errorResponse
// @Router /appointments/{id} [get]
// @Security OAuth2Application
func (h *appointmentHandler) GetByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid id provided")
			return
		}
		response, err := h.s.GetByID(id)
		if err != nil {
			web.BadResponse(ctx, http.StatusNotFound, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}

// GetAllByIdentityNumber - get all appointments by patient identity doc
// @BasePath /api/v1
// GetAllByIdentityNumber godoc
// @Summary Get all appointments by patient identity doc
// @Schemes
// @Description get all appointments by patient identity doc
// @Tags Appointments
// @Accept json
// @Produce json
// @Param identity_number path int true "Patient Doc Number"
// @Success 200 {object} []domain.AppointmentDTO
// @Failure 400 {object} web.errorResponse
// @Failure 401 {object} web.errorResponse
// @Failure 404 {object} web.errorResponse
// @Router /appointments/patient/{identity_number} [get]
// @Security OAuth2Application
func (h *appointmentHandler) GetAllByIdentityNumber() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("identity_number")
		response, err := h.s.GetAllByIdentityNumber(idParam)
		if err != nil {
			web.BadResponse(ctx, http.StatusNotFound, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}

// GetAllByLicenseNumber - get all appointments by dentist license doc
// @BasePath /api/v1
// GetAllByLicenseNumber godoc
// @Summary Get all appointments by dentist license doc
// @Schemes
// @Description Get all appointments by dentist license doc
// @Tags Appointments
// @Accept json
// @Produce json
// @Param license_number path int true "Dentist License Number"
// @Success 200 {object} []domain.AppointmentDTO
// @Failure 400 {object} web.errorResponse
// @Failure 401 {object} web.errorResponse
// @Failure 404 {object} web.errorResponse
// @Router /appointments/dentist/{license_number} [get]
// @Security OAuth2Application
func (h *appointmentHandler) GetAllByLicenseNumber() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("license_number")
		response, err := h.s.GetAllByLicenseNumber(idParam)
		if err != nil {
			web.BadResponse(ctx, http.StatusNotFound, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}

// Post - save a new appointment
// @BasePath /api/v1
// Post godoc
// @Summary Create a new appointment
// @Schemes
// @Description Create a new appointment by request body
// @Tags Appointments
// @Accept json
// @Produce json
// @Param body body domain.Appointment true "Body"
// @Success 201 {object} domain.AppointmentDTO
// @Failure 400 {object} web.errorResponse
// @Failure 401 {object} web.errorResponse
// @Router /appointments [post]
// @Security OAuth2Application
func (h *appointmentHandler) Post() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var appointment domain.Appointment
		err := ctx.ShouldBindJSON(&appointment)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid appointment data, please verify field(s): "+err.Error())
			return
		}

		isValid, err := isEmptyAppointment(&appointment)
		if !isValid {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}
		response, err := h.s.Create(appointment)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusCreated, response)
	}
}

// Put - update an entire appointment by ID
// @BasePath /api/v1
// PutAppointment godoc
// @Summary Update a entire appointment by ID
// @Schemes
// @Description Update a entire appointment by ID and request body
// @Tags Appointments
// @Accept json
// @Produce json
// @Param id path int true "Appointment ID"
// @Param body body domain.Appointment true "Body"
// @Success 200 {object} domain.AppointmentDTO
// @Failure 400 {object} web.errorResponse
// @Failure 401 {object} web.errorResponse
// @Failure 404 {object} web.errorResponse
// @Failure 409 {object} web.errorResponse
// @Router /appointments/{id} [put]
// @Security OAuth2Application
func (h *appointmentHandler) Put() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid id provided")
			return
		}

		var appointment domain.Appointment
		err = ctx.ShouldBindJSON(&appointment)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid appointment data, verify the fields and try again")
			return
		}

		isValid, err := isEmptyAppointment(&appointment)
		if !isValid {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}
		response, err := h.s.Update(id, appointment)
		if err != nil {
			web.BadResponse(ctx, http.StatusNotFound, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}

// Patch - update fields from an appointment by ID
// @BasePath /api/v1
// PatchAppointment godoc
// @Summary Update fields from an appointment by ID
// @Schemes
// @Description Update fields from an appointment by ID and fields in request body
// @Tags Appointments
// @Accept json
// @Produce json
// @Param id path int true "Appointment ID"
// @Param body body domain.Appointment true "Body"
// @Success 200 {object} domain.AppointmentDTO
// @Failure 400 {object} web.errorResponse
// @Failure 401 {object} web.errorResponse
// @Failure 404 {object} web.errorResponse
// @Router /appointments/{id} [patch]
// @Security OAuth2Application
func (h *appointmentHandler) Patch() gin.HandlerFunc {
	type Request struct {
		Description string `json:"description,omitempty"`
		DateAndTime string `json:"dateAndTime,omitempty"`
		DentistCRO  string `json:"dentistCRO,omitempty"`
		PatientRG   string `json:"patientRG,omitempty"`
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
		update := domain.Appointment{
			Description: r.Description,
			DateAndTime: r.DateAndTime,
			DentistCRO:  r.DentistCRO,
			PatientRG:   r.PatientRG,
		}
		if update.DateAndTime != "" {
			if !validateDateTime(update.DateAndTime) {
				web.BadResponse(ctx, http.StatusBadRequest, "error", "please the appointment must be in format: 30/01/2023 23:59")
				return
			}
		}
		response, err := h.s.Update(id, update)
		if err != nil {
			web.BadResponse(ctx, http.StatusNotFound, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}

// Delete - delete an appointment
// @BasePath /api/v1
// DeleteAppointment godoc
// @Summary Delete an appointment by ID
// @Schemes
// @Description Delete an appointment by ID
// @Tags Appointments
// @Accept json
// @Produce json
// @Param id path int true "Appointment ID"
// @Success 200 {object} web.errorResponse
// @Failure 400 {object} web.errorResponse
// @Failure 401 {object} web.errorResponse
// @Failure 404 {object} web.errorResponse
// @Router /appointments/{id} [delete]
// @Security OAuth2Application
func (h *appointmentHandler) Delete() gin.HandlerFunc {
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
		web.DeleteResponse(ctx, http.StatusOK, "appointment removed")
	}
}

// Aux functions bellow->

func isEmptyAppointment(appointment *domain.Appointment) (bool, error) {
	dateTimeParsed, err := time.Parse("02/01/2006 15:04", appointment.DateAndTime)
	if err != nil {
		return false, err
	}
	switch {
	case appointment.Description == "" || appointment.DentistCRO == "" || appointment.DateAndTime == "" || appointment.PatientRG == "":
		return false, errors.New("fields can't be empty")
	case !validateDateTime(appointment.DateAndTime):
		return false, errors.New("please the appointment must be in format: 30/01/2023 23:59")
	case dateTimeParsed.Local().Add(time.Hour * 3).Before(time.Now().Add(time.Hour)):
		return false, errors.New("the appointment must be in +1 hour from now")
	}
	return true, nil
}

func validateDateTime(dateTime string) bool {
	datesInit := strings.Split(dateTime, " ")
	if len(datesInit) != 2 {
		log.Printf("invalid time, must be in format: 30/01/2023 23:59")
		return false
	}
	breakDate := strings.Split(datesInit[0], "/")
	if len(breakDate) != 3 {
		log.Println("invalid time, must be in format: 30/01/2023 23:59 or 30/01/2023 23:59:59")
		return false
	}
	breakTime := strings.Split(datesInit[1], ":")
	var listDate []int
	var listTime []int

	for _, date := range breakDate {
		number, err := strconv.Atoi(date)
		if err != nil {
			return false
		}
		listDate = append(listDate, number)
	}
	condition := (listDate[0] < 1 || listDate[0] > 31) && (listDate[1] < 1 || listDate[1] > 12) && (listDate[2] < 1 || listDate[2] > 9999)
	if condition {
		log.Println("invalid time, must be between: 1 and 31/12/2023 23:59")
		return false
	}

	for _, t := range breakTime {
		clock, err := strconv.Atoi(t)
		if err != nil {
			log.Println("invalid time, must be in format 23:59 (hours and minutes)")
			return false
		}
		listTime = append(listTime, clock)
	}

	if len(listTime) == 2 {
		condition = (listTime[0] < 0 || listTime[0] > 23) && (listTime[1] < 0 || listTime[1] > 59)
		if condition {
			log.Println("invalid time, must be between: 00:00 and 23:59")
			return false
		}
	}

	if len(listTime) == 3 {
		condition = (listTime[0] < 0 || listTime[0] > 23) && (listTime[1] < 0 || listTime[1] > 59)
		if condition {
			log.Println("invalid time, must be between: 00:00 and 23:59")
			return false
		}
	}
	return true
}
