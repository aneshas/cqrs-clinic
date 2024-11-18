package api

import (
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

type patient struct {
	ID            string  `bson:"patient_id" json:"id"`
	Name          string  `bson:"patient_name" json:"name"`
	WardNumber    string  `bson:"ward_number" json:"ward_number"`
	Age           int     `bson:"patient_age" json:"age"`
	Status        string  `bson:"status" json:"status"`
	DischargeNote *string `bson:"discharge_note" json:"discharge_note,omitempty"`
} // @name Patient

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

	results := make([]patient, 0)

	err = cursor.All(ctx, &results)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, results)
}
