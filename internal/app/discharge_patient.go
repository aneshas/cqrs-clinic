package app

import (
	"context"
	"github.com/aneshas/clinic/internal/patient"
	"github.com/aneshas/eventstore/aggregate"
)

// DischargePatientFunc represents discharge patient use case
type DischargePatientFunc func(ctx context.Context, patientID string) error

// NewDischargePatient creates a new discharge patient use case
func NewDischargePatient(store *aggregate.Store[*patient.Patient]) DischargePatientFunc {
	exec := aggregate.NewExecutor[*patient.Patient](store)

	return func(ctx context.Context, patientID string) error {
		p := patient.New(patientID)

		return exec(ctx, p, func(ctx context.Context) error {
			p.Discharge()

			return nil
		})
	}
}
