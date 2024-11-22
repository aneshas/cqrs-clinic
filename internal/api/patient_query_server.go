package api

import (
	"github.com/aneshas/clinic/internal/app"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

// RegisterPatientQueryServer registers patient query server
func RegisterPatientQueryServer(
	e *echo.Echo,
	client *mongo.Client,
) {
	s := &PatientQueryServer{
		client: client,
		db:     client.Database("clinic"),
	}

	patients := e.Group("/patients")

	patients.GET("", s.patients)
}

// PatientQueryServer represents a patient query server
type PatientQueryServer struct {
	client *mongo.Client
	db     *mongo.Database
}

// @Summary Get patient roster
// @Description Get all patients.
// @Tags Patients
// @Accept json
// @Produce json
// @Success 200 {array} Patient "All patients"
// @Failure 400 {object} HttpError
// @Failure 500 {object} HttpError
// @Router /patients [get]
func (s *PatientQueryServer) patients(c echo.Context) error {
	ctx := c.Request().Context()

	cursor, err := s.db.Collection("patient_roster").Find(ctx, bson.D{})
	if err != nil {
		return err
	}

	results := make([]app.Patient, 0)

	err = cursor.All(ctx, &results)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, results)
}
