package api

import (
	"github.com/aneshas/clinic/internal/app"
	"github.com/aneshas/eventstore"
	"github.com/aneshas/eventstore/ambar"
	"github.com/aneshas/eventstore/ambar/echoambar"
	"github.com/labstack/echo/v4"
)

// RegisterProjectionServer registers projection server
func RegisterProjectionServer(
	e *echo.Echo,
	patientRosterProjection ambar.Projection,
) {
	projections := e.Group("/projections")

	projections.POST(
		"/patient_roster",
		handler(app.PatientRosterSubscriptions...)(patientRosterProjection),
	)
}

func handler(events ...any) func(projection ambar.Projection) echo.HandlerFunc {
	return echoambar.Wrap(
		ambar.New(eventstore.NewJSONEncoder(events...)),
	)
}
