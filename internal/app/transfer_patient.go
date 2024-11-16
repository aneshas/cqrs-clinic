package app

import (
	"context"
	"github.com/aneshas/clinic/internal/patient"
	"github.com/aneshas/eventstore/aggregate"
)

// PatientToTransfer represents the command to transfer a patient to a new ward
type PatientToTransfer struct {
	PatientID     string
	NewWardNumber string
}

// TransferPatientFunc represents transfer patient use case
type TransferPatientFunc func(ctx context.Context, cmd PatientToTransfer) error

// NewTransferPatient creates a new transfer patient use case
func NewTransferPatient(store *aggregate.Store[*patient.Patient]) TransferPatientFunc {
	exec := aggregate.NewExecutor[*patient.Patient](store)

	return func(ctx context.Context, cmd PatientToTransfer) error {
		p := patient.New(cmd.PatientID)

		return exec(ctx, p, func(ctx context.Context) error {
			return p.TransferTo(cmd.NewWardNumber)
		})
	}
}
