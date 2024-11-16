package patient

// Events is a list of all patient domain event instances
var Events = []any{
	Admitted{},
	Transferred{},
	Discharged{},
}

// Admitted event signals that a patient has been admitted
type Admitted struct {
	PatientID   string
	PatientName string
	WardNumber  string
	PatientAge  int
}

// Transferred event signals that a patient has been transferred to a new ward
type Transferred struct {
	PatientID     string
	NewWardNumber string
}

// Discharged event signals that a patient has been discharged
type Discharged struct {
	PatientID     string
	DischargeNote string
}
