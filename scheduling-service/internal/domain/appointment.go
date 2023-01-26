package domain

type Appointment struct {
	Id          int    `json:"id"`
	Description string `json:"description" binding:"required"`
	DateAndTime string `json:"dateAndTime" binding:"required"`
	DentistCRO  string `json:"dentistCRO" binding:"required"`
	PatientRG   string `json:"patientRG" binding:"required"`
}
