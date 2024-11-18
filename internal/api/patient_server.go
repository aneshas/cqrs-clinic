package api

import (
	"github.com/aneshas/clinic/internal/app"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"

	_ "github.com/aneshas/clinic/docs"
)

// RegisterPatientServer registers patient server
func RegisterPatientServer(
	e *echo.Echo,
	admitPatient app.AdmitPatientFunc,
	transferPatient app.TransferPatientFunc,
	dischargePatient app.DischargePatientFunc,
) {
	s := &PatientServer{
		admitPatient:     admitPatient,
		transferPatient:  transferPatient,
		dischargePatient: dischargePatient,
	}

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	patients := e.Group("/patients")

	patients.POST("/admit", s.admitPatientHandler)
	patients.POST("/:id/transfer", s.transferPatientHandler)
	patients.POST("/:id/discharge", s.dischargePatientHandler)
}

// PatientServer represents a patient http server
type PatientServer struct {
	admitPatient     app.AdmitPatientFunc
	transferPatient  app.TransferPatientFunc
	dischargePatient app.DischargePatientFunc
}

// Patient to be admitted
// @Description Patient to be admitted
type admitPatientRequest struct {
	PatientName string `json:"patient_name"`
	WardNumber  string `json:"ward_number"`
	PatientAge  int    `json:"patient_age"`
} // @name AdmitPatientRequest

type admitPatientResponse struct {
	PatientID string `json:"patient_id"`
} // @name AdmitPatientResponse

// @Summary Admit a new patient
// @Description admits a new patient.
// @Tags Patients
// @Accept json
// @Produce json
// @Param AdmitPatientRequest body AdmitPatientRequest true "Patient to admit"
// @Success 200 {object} AdmitPatientResponse "Admitted patient"
// @Failure 400 {object} HttpError
// @Failure 500 {object} HttpError
// @Router /patients/admit [post]
func (s *PatientServer) admitPatientHandler(c echo.Context) error {
	var req admitPatientRequest

	err := c.Bind(&req)
	if err != nil {
		return err
	}

	admittedPatient, err := s.admitPatient(
		c.Request().Context(),
		app.PatientToAdmit{
			PatientName: req.PatientName,
			WardNumber:  req.WardNumber,
			PatientAge:  req.PatientAge,
		},
	)
	if err != nil {
		return err
	}

	return c.JSON(
		http.StatusOK,
		admitPatientResponse{
			PatientID: admittedPatient.PatientID,
		},
	)
}

type transferPatientRequest struct {
	NewWardNumber string `json:"new_ward_number"`
} // @name TransferPatientRequest

// @Summary Transfer a patient
// @Description transfers patient to a new ward.
// @Tags Patients
// @Accept json
// @Produce json
// @Param id path string true "Patient Id"
// @Param TransferPatientRequest body TransferPatientRequest true "Transfer request"
// @Success 204
// @Failure 400 {object} HttpError
// @Failure 500 {object} HttpError
// @Router /patients/{id}/transfer [post]
func (s *PatientServer) transferPatientHandler(c echo.Context) error {
	var req transferPatientRequest

	err := c.Bind(&req)
	if err != nil {
		return err
	}

	err = s.transferPatient(
		c.Request().Context(),
		app.PatientToTransfer{
			PatientID:     c.Param("id"),
			NewWardNumber: req.NewWardNumber,
		},
	)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

type dischargePatientRequest struct {
	DischargeNote string `json:"discharge_note"`
} // @name DischargePatientRequest

// @Summary Discharge a patient
// @Description discharges a patient.
// @Tags Patients
// @Accept json
// @Produce json
// @Param id path string true "Patient Id"
// @Param DischargePatientRequest body DischargePatientRequest true "Discharge request"
// @Success 204
// @Failure 400 {object} HttpError
// @Failure 500 {object} HttpError
// @Router /patients/{id}/discharge [post]
func (s *PatientServer) dischargePatientHandler(c echo.Context) error {
	var req dischargePatientRequest

	err := c.Bind(&req)
	if err != nil {
		return err
	}

	err = s.dischargePatient(
		c.Request().Context(),
		app.PatientToDischarge{
			PatientID:     c.Param("id"),
			DischargeNote: req.DischargeNote,
		},
	)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
