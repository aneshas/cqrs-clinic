package app

import (
	"context"
	"github.com/aneshas/clinic/internal/patient"
	"github.com/aneshas/eventstore/aggregate"
)

// PatientToAdmit represents the command to admit a patient
type PatientToAdmit struct {
	PatientID   string
	PatientName string
	WardNumber  string
	PatientAge  int
}

// AdmittedPatient represents the admitted patient
type AdmittedPatient struct {
	PatientID string
}

// AdmitPatientFunc represents admit patient use case
type AdmitPatientFunc func(ctx context.Context, cmd PatientToAdmit) (*AdmittedPatient, error)

// NewAdmitPatient creates a new admit patient use case
func NewAdmitPatient(store *aggregate.Store[*patient.Patient]) AdmitPatientFunc {
	return func(ctx context.Context, cmd PatientToAdmit) (*AdmittedPatient, error) {
		p, err := patient.NewForAdmission(patient.NewID(), cmd.PatientName, cmd.WardNumber, cmd.PatientAge)
		if err != nil {
			return nil, err
		}

		err = store.Save(ctx, p)
		if err != nil {
			return nil, err
		}

		return &AdmittedPatient{p.StringID()}, nil
	}
}
