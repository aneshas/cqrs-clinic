package app

import (
	"fmt"
	"github.com/aneshas/clinic/internal/patient"
	"github.com/aneshas/eventstore"
	"github.com/aneshas/eventstore/ambar"
	"net/http"
)

// PatientRosterSubscriptions represents patient roster event subscriptions
var PatientRosterSubscriptions = []any{
	patient.Admitted{},
	patient.Transferred{},
	patient.Discharged{},
}

// NewPatientRosterProjection creates a new patient roster projection
func NewPatientRosterProjection() ambar.Projection {
	return func(_ *http.Request, event eventstore.StoredEvent) error {
		switch event.Event.(type) {
		case patient.Admitted:
			evt := event.Event.(patient.Admitted)

			fmt.Printf("Patient: #%s | Admitted to ward: <%s>\n", evt.PatientID, evt.WardNumber)

		default:
			fmt.Println("not interested in this event")
		}

		return nil
	}
}
