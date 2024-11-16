package app

import (
	"context"
	"github.com/aneshas/clinic/internal/patient"
	"github.com/aneshas/eventstore/aggregate"
)

// PatientToDischarge represents a patient to discharge
type PatientToDischarge struct {
	PatientID     string
	DischargeNote string
}

// DischargePatientFunc represents discharge patient use case
type DischargePatientFunc func(ctx context.Context, cmd PatientToDischarge) error

// NewDischargePatient creates a new discharge patient use case
func NewDischargePatient(store *aggregate.Store[*patient.Patient]) DischargePatientFunc {
	exec := aggregate.NewExecutor[*patient.Patient](store)

	return func(ctx context.Context, cmd PatientToDischarge) error {
		p := patient.New(cmd.PatientID)

		return exec(ctx, p, func(ctx context.Context) error {
			p.Discharge(cmd.DischargeNote)

			return nil
		})
	}
}
