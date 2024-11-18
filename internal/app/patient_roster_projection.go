package app

import (
	"context"
	"fmt"
	"github.com/aneshas/clinic/internal/patient"
	"github.com/aneshas/eventstore"
	"github.com/aneshas/eventstore/ambar"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

// PatientRosterSubscriptions represents patient roster event subscriptions
var PatientRosterSubscriptions = []any{
	patient.Admitted{},
	patient.Transferred{},
	patient.Discharged{},
}

// NewPatientRosterProjection creates a new patient roster projection
func NewPatientRosterProjection(client *mongo.Client) ambar.Projection {
	ctx := context.Background()
	coll := client.Database("clinic").Collection("patient_roster")

	return func(_ *http.Request, event eventstore.StoredEvent) error {
		switch event.Event.(type) {
		case patient.Admitted:
			evt := event.Event.(patient.Admitted)

			fmt.Printf("Patient: #%s | Admitted to ward: <%s>\n", evt.PatientID, evt.WardNumber)

			_, err := coll.InsertOne(ctx, bson.D{
				{"patient_id", evt.PatientID},
				{"patient_name", evt.PatientName},
				{"ward_number", evt.WardNumber},
				{"patient_age", evt.PatientAge},
				{"status", "admitted"},
			})

			return err

		case patient.Transferred:
			evt := event.Event.(patient.Transferred)

			fmt.Printf("Patient: #%s | Transferred to ward: <%s>\n", evt.PatientID, evt.NewWardNumber)

			_, err := coll.UpdateOne(
				ctx,
				bson.D{
					{"patient_id", evt.PatientID},
				},
				bson.M{
					"$set": bson.D{
						{"ward_number", evt.NewWardNumber},
					},
				},
			)

			return err

		case patient.Discharged:
			evt := event.Event.(patient.Discharged)

			fmt.Printf("Patient: #%s | Discharged\n", evt.PatientID)

			_, err := coll.UpdateOne(
				ctx,
				bson.D{
					{"patient_id", evt.PatientID},
				},
				bson.M{
					"$set": bson.D{
						{"status", "discharged"},
						{"discharge_note", evt.DischargeNote},
					},
				},
			)

			return err

		default:
			fmt.Println("not interested in this event")
		}

		return nil
	}
}
