package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ronilsonalves/GoLang-in-a-spring-cloud-architecture/scheduling-service/internal/dentist"
	"github.com/ronilsonalves/GoLang-in-a-spring-cloud-architecture/scheduling-service/internal/domain"
	"github.com/ronilsonalves/GoLang-in-a-spring-cloud-architecture/scheduling-service/pkg/web"
	"net/http"
	"strconv"
)

type dentistHandler struct {
	s dentist.Service
}

func NewDentistHandler(s dentist.Service) *dentistHandler {
	return &dentistHandler{
		s: s,
	}
}

// GetAll - get all dentists from db.
// @BasePath /api/v1
// GetAllDentists godoc
// @Summary List all dentists
// @Schemes
// @Description get all dentists from db.
// @Tags Dentists
// @Accept json
// @Produce json
// @Success 200 {object} []domain.Dentist
// @Failure 400 {object} web.errorResponse
// @Failure 401 {object} web.errorResponse
// @Router /dentists [get]
// @Security SECRET_TOKEN
func (h *dentistHandler) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		response, err := h.s.GetAll()
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}

// GetByID - get a dentist by an ID
// @BasePath /api/v1
// GetDentistByID godoc
// @Summary Get a dentist by an ID
// @Schemes
// @Description get a dentist by a provided ID.
// @Tags Dentists
// @Accept json
// @Produce json
// @Param id path int true "Dentist ID"
// @Success 200 {object} domain.Dentist
// @Failure 400 {object} web.errorResponse
// @Failure 401 {object} web.errorResponse
// @Failure 404 {object} web.errorResponse
// @Router /dentists/{id} [get]
// @Security SECRET_TOKEN
func (h *dentistHandler) GetByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid id provided")
			return
		}

		response, err := h.s.GetByID(id)
		if err != nil {
			web.BadResponse(ctx, http.StatusNotFound, "error", "dentist not found")
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}

// Post - save a new dentist into db
// @BasePath /api/v1
// PostDentist godoc
// @Summary Create a new dentist
// @Schemes
// @Description Create a new dentist by request body.
// @Tags Dentists
// @Accept json
// @Produce json
// @Param body body domain.Dentist true "Body"
// @Success 201 {object} domain.Dentist
// @Failure 400 {object} web.errorResponse
// @Failure 401 {object} web.errorResponse
// @Router /dentists [post]
// @Security SECRET_TOKEN
func (h *dentistHandler) Post() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var dentist domain.Dentist
		err := ctx.ShouldBindJSON(&dentist)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid dentist")
			return
		}

		isValid, err := isEmptyDentist(&dentist)
		if !isValid {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}

		response, err := h.s.Create(dentist)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusCreated, response)
	}
}

// Put - update an entire dentist
// @BasePath /api/v1
// PutDentist godoc
// @Summary Update an entire dentist by ID
// @Schemes
// @Description Update an entire dentist by ID.
// @Tags Dentists
// @Accept json
// @Produce json
// @Param id path int true "Dentist ID"
// @Param body body domain.Dentist true "Body"
// @Success 200 {object} domain.Dentist
// @Failure 400 {object} web.errorResponse
// @Failure 401 {object} web.errorResponse
// @Failure 409 {object} web.errorResponse
// @Router /dentists/{id} [put]
// @Security SECRET_TOKEN
func (h *dentistHandler) Put() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid id")
			return
		}
		var dentist domain.Dentist
		err = ctx.ShouldBindJSON(&dentist)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid dentist data")
		}

		isValid, err := isEmptyDentist(&dentist)
		if !isValid {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}

		response, err := h.s.Update(id, dentist)
		if err != nil {
			web.BadResponse(ctx, http.StatusConflict, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}

// Patch - update fields from a dentist
// @BasePath /api/v1
// PatchDentist godoc
// @Summary Update fields from a dentist
// @Schemes
// @Description Update fields from a dentist
// @Tags Dentists
// @Accept json
// @Produce json
// @Param id path int true "Dentist ID"
// @Param body body domain.Dentist true "Body"
// @Success 200 {object} domain.Dentist
// @Failure 400 {object} web.errorResponse
// @Failure 401 {object} web.errorResponse
// @Router /dentists/{id} [patch]
// @Security SECRET_TOKEN
func (h *dentistHandler) Patch() gin.HandlerFunc {
	type Request struct {
		Surname       string `json:"surname,omitempty"`
		Name          string `json:"name,omitempty"`
		LicenseNumber string `json:"license_number,omitempty"`
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
		update := domain.Dentist{
			LastName: r.Surname,
			Name:     r.Name,
			CRO:      r.LicenseNumber,
		}

		updated, err := h.s.Update(id, update)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, updated)
	}
}

// Delete - delete a dentist
// @BasePath /api/v1
// DeleteDentist godoc
// @Summary Delete a product by ID
// @Schemes
// @Description Delete a product by ID
// @Tags Dentists
// @Accept json
// @Produce json
// @Param id path int true "Dentist ID"
// @Success 200 {object} web.response
// @Failure 400 {object} web.errorResponse
// @Failure 404 {object} web.errorResponse
// @Failure 401 {object} web.errorResponse
// @Router /dentists/{id} [delete]
// @Security SECRET_TOKEN
func (h *dentistHandler) Delete() gin.HandlerFunc {
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
		web.DeleteResponse(ctx, http.StatusOK, "dentist deleted")
	}
}

// Aux functions bellow->

func isEmptyDentist(dentist *domain.Dentist) (bool, error) {
	switch {
	case dentist.LastName == "" || dentist.Name == "" || dentist.CRO == "":
		return false, errors.New("fields can't be empty")
	}
	return true, nil
}
