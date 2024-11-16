package api

import (
	"github.com/aneshas/clinic/internal/app"
	"github.com/labstack/echo/v4"
	"net/http"
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

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	patients := e.Group("/patients")

	patients.POST("", s.admitPatientHandler)
	patients.POST("/:id/transfer", s.transferPatientHandler)
	patients.GET("/:id/discharge", s.dischargePatientHandler)
}

// PatientServer represents a patient http server
type PatientServer struct {
	admitPatient     app.AdmitPatientFunc
	transferPatient  app.TransferPatientFunc
	dischargePatient app.DischargePatientFunc
}

type admitPatientRequest struct {
	PatientName string `json:"patient_name"`
	WardNumber  string `json:"ward_number"`
	PatientAge  int    `json:"patient_age"`
}

type admitPatientResponse struct {
	PatientID string `json:"patient_id"`
}

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
}

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

func (s *PatientServer) dischargePatientHandler(c echo.Context) error {
	err := s.dischargePatient(c.Request().Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
