package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"transport-status/pkg"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type VehicleData struct {
	DB *sql.DB
}

const defaultReisaiFileName = "reisai.txt"

func CreateTransportDatabase(db *sql.DB) (*VehicleData, error) {

	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS transport(VehicleType CHAR, Route CHAR, Schedule INTEGER, Shift INTEGER, BusNumber CHAR PRIMARY KEY, LowGrind INTEGER, TripsStart INTEGER, TripsEnd INTEGER, DirectionID INTEGER, DirectionType CHAR, DirectionName CHAR )")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create table transport")
	}
	_, err = statement.Exec()
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to execute create table statement")
	}
	return &VehicleData{DB: db}, err
}

func (vehicledata *VehicleData) UpdateTransportDatabase() error {

	file, err := os.Open(defaultReisaiFileName)
	if err != nil {
		return errors.Wrap(err, "failed to open Reisai file")
	}
	defer file.Close()
	dataSlice, err := pkg.ParseData(file)
	if err != nil {
		return errors.Wrap(err, "failed to parse data")
	}

	statement, err := vehicledata.DB.Prepare("INSERT INTO transport (VehicleType,Route, Schedule, Shift, BusNumber, LowGrind, TripsStart, TripsEnd, DirectionID, DirectionType, DirectionName) VALUES (?,?,?,?,?,?,?,?,?,?,?) ON CONFLICT(BusNumber) DO NOTHING")
	if err != nil {
		return errors.Wrap(err, "failed to prepare database")
	}
	for _, data := range dataSlice {
		_, err = statement.Exec(data.VehicleType, data.Route, data.Schedule, data.Shift, data.BusNumber, data.LowGrind, data.TripsStart, data.TripsEnd, data.DirectionId, data.DirectionType, data.DirectionName)
		if err != nil {
			return errors.Wrap(err, "failed to execute data insert statement")
		}
	}
	return nil
}

func (vehicleData *VehicleData) GetTransports() ([]Transport, error) {

	var transport Transport

	transports := []Transport{}
	rows, err := vehicleData.DB.Query("SELECT  * FROM transport")
	if err != nil {
		return nil, errors.Wrap(err, "failed to Query 	")
	}
	for rows.Next() {
		rows.Scan(&transport.VehicleType, &transport.Route, &transport.Schedule, &transport.Shift, &transport.BusNumber, &transport.LowGrind, &transport.TripsStart, &transport.TripsEnd, &transport.DirectionId, &transport.DirectionType, &transport.DirectionType)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan data from database")
		}
		transports = append(transports, transport)
	}

	return transports, nil
}

type Route struct {
	BusNumber string
	Route     string
}

func (vehicleData *VehicleData) GetRoute(BusNumber string) ([]Route, error) {
	var transport Route
	transports := []Route{}
	rows, err := vehicleData.DB.Query("SELECT BusNumber ,Route FROM transport WHERE BusNumber =?", BusNumber)
	if err != nil {
		fmt.Println("error")
		fmt.Println(err)
		return nil, errors.Wrap(err, "failed to Query 	")
	}
	for rows.Next() {
		rows.Scan(&transport.BusNumber, &transport.Route)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan data from database")
		}
		transports = append(transports, transport)
	}

	return transports, nil
}

func ConnectionToTransportDatabase() *sql.DB {

	transportDatabase, err := sql.Open("sqlite3", "../../transport.db")
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to open transport database"))
	}

	return transportDatabase
}
