package patient

import (
	"context"
	"fmt"
	"github.com/aneshas/clinic/internal/errs"
	"github.com/aneshas/eventstore/aggregate"
)

var (
	// ErrInvalidPatientName represents an error when patient name is invalid
	ErrInvalidPatientName = errs.E(fmt.Errorf("patient name is required"))

	// ErrInvalidWardNumber represents an error when ward number is invalid
	ErrInvalidWardNumber = errs.E(fmt.Errorf("ward number is required"))

	// ErrInvalidPatientAge represents an error when patient age is invalid
	ErrInvalidPatientAge = errs.E(fmt.Errorf("patient age is invalid"))

	// ErrPatientDischarged represents an error when patient is already discharged
	ErrPatientDischarged = errs.E(fmt.Errorf("patient is already discharged"))

	// ErrSameWardNumber represents an error when transferring to the same ward
	ErrSameWardNumber = errs.E(fmt.Errorf("cannot transfer to the same ward"))
)

// Store represents a patient store
type Store interface {
	ByID(ctx context.Context, id string, p *Patient) error
	Save(ctx context.Context, p *Patient) error
}

// New instantiates a new Patient with an ID
func New(id string) *Patient {
	var p Patient

	p.ID = ParseID(id)

	return &p
}

// NewForAdmission instantiates a new Patient
func NewForAdmission(id ID, name, ward string, age int) (*Patient, error) {
	if name == "" {
		return nil, ErrInvalidPatientName
	}

	if ward == "" {
		return nil, ErrInvalidWardNumber
	}

	if age < 0 {
		return nil, ErrInvalidPatientAge
	}

	var p Patient

	p.Rehydrate(&p)

	p.Apply(
		Admitted{
			PatientID:   id.String(),
			PatientName: name,
			WardNumber:  ward,
			PatientAge:  age,
		},
	)

	return &p, nil
}

// Patient represents a patient aggregate
type Patient struct {
	aggregate.Root[ID]

	isDischarged bool
	ward         string
}

// TransferTo transfers a patient to a new ward
func (p *Patient) TransferTo(newWard string) error {
	if p.isDischarged {
		return ErrPatientDischarged
	}

	if p.ward == newWard {
		return ErrSameWardNumber
	}

	p.Apply(
		Transferred{
			PatientID:     p.ID.String(),
			NewWardNumber: newWard,
		},
	)

	return nil
}

// Discharge discharges a patient
func (p *Patient) Discharge(note string) {
	if p.isDischarged {
		return
	}

	p.Apply(
		Discharged{
			PatientID:     p.ID.String(),
			DischargeNote: note,
		},
	)
}

// OnAdmitted event handler
func (p *Patient) OnAdmitted(e Admitted) {
	p.ID = ParseID(e.PatientID)
	p.ward = e.WardNumber
}

// OnTransferred event handler
func (p *Patient) OnTransferred(e Transferred) {
	p.ward = e.NewWardNumber
}

// OnDischarged event handler
func (p *Patient) OnDischarged(_ Discharged) {
	p.isDischarged = true
}
