package main

import (
	"errors"
	"flag"
	"github.com/aneshas/clinic/internal/api"
	"github.com/aneshas/clinic/internal/app"
	"github.com/aneshas/clinic/internal/errs"
	"github.com/aneshas/clinic/internal/patient"
	"github.com/aneshas/eventstore"
	"github.com/aneshas/eventstore/aggregate"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"os"
)

var pg = flag.Bool("pg", false, "Run with postgres db (set DSN env to pg connection string)")

// @title Event-Sourced Clinic Example API
// @version 1.0
// @description This is an Event-Sourced Clinic example.
func main() {
	flag.Parse()

	db := eventstore.WithSQLiteDB("clinic.db")

	if *pg {
		db = eventstore.WithPostgresDB(os.Getenv("DSN"))
	}

	eventStore, err := eventstore.New(
		eventstore.NewJSONEncoder(patient.Events...),
		db,
	)
	checkErr(err)

	if *pg {
		m, err := migrate.New(os.Getenv("MIGRATION_SOURCE"), os.Getenv("MIGRATION_DSN"))
		checkErr(err)

		checkErr(m.Up())
	}

	e := echo.New()

	e.HTTPErrorHandler = errs.ErrorHandler

	e.Use(
		middleware.Logger(),
		middleware.Recover(),
	)

	patientStore := aggregate.NewStore[*patient.Patient](eventStore)

	api.RegisterPatientServer(
		e,
		app.NewAdmitPatient(patientStore),
		app.NewTransferPatient(patientStore),
		app.NewDischargePatient(patientStore),
	)

	api.RegisterProjectionServer(
		e,
		app.NewPatientRosterProjection(),
	)

	log.Fatal(e.Start(":8080"))
}

func checkErr(err error) {
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return
		}

		log.Fatal(err)
	}
}
