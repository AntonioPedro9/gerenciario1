package handlers

import (
	"net/http"
	"server/models"
	"server/services"
	"server/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type AppointmentHandler struct {
	appointmentService *services.AppointmentService
}

func NewAppointmentHandler(appointmentService *services.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{appointmentService}
}

func (ah *AppointmentHandler) CreateAppointment(c *gin.Context) {
	var appointment models.CreateAppointmentRequest
	if err := c.ShouldBindJSON(&appointment); err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind JSON request"})
		return
	}

	if err := ah.appointmentService.CreateAppointment(&appointment); err != nil {
		log.Error(err)

		customError, ok := err.(*utils.CustomError)
		if !ok {
			c.JSON(http.StatusInternalServerError, nil)
			return
		}

		c.JSON(customError.StatusCode, gin.H{"error": customError.Message})
		return
	}

	c.JSON(http.StatusCreated, nil)
}

func (ah *AppointmentHandler) ListAppointments(c *gin.Context) {
	id := c.Param("userID")

	userID, err := uuid.Parse(id)
	if err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No token provided"})
		return
	}

	authUserID, err := utils.GetIDFromToken(tokenString)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	appointments, err := ah.appointmentService.ListAppointments(userID, authUserID)
	if err != nil {
		log.Error(err)

		customError, ok := err.(*utils.CustomError)
		if !ok {
			c.JSON(http.StatusInternalServerError, nil)
			return
		}

		c.JSON(customError.StatusCode, gin.H{"error": customError.Message})
		return
	}

	c.JSON(http.StatusOK, appointments)
}

func (ah *AppointmentHandler) UpdateAppointment(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No token provided"})
		return
	}

	authUserID, err := utils.GetIDFromToken(tokenString)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var appointment models.UpdateAppointmentRequest
	if err := c.ShouldBindJSON(&appointment); err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind JSON request"})
		return
	}

	if err := ah.appointmentService.UpdateAppointment(&appointment, authUserID); err != nil {
		log.Error(err)

		customError, ok := err.(*utils.CustomError)
		if !ok {
			c.JSON(http.StatusInternalServerError, nil)
			return
		}

		c.JSON(customError.StatusCode, gin.H{"error": customError.Message})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (ah *AppointmentHandler) DeleteAppointment(c *gin.Context) {
	id := c.Param("appointmentID")

	parsedID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}
	appointmentID := uint(parsedID)

	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No token provided"})
		return
	}

	authUserID, err := utils.GetIDFromToken(tokenString)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := ah.appointmentService.DeleteAppointment(appointmentID, authUserID); err != nil {
		log.Error(err)

		customError, ok := err.(*utils.CustomError)
		if !ok {
			c.JSON(http.StatusInternalServerError, nil)
			return
		}

		c.JSON(customError.StatusCode, gin.H{"error": customError.Message})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}