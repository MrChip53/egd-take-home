package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"os"
	"take-home/pkg/data"
)

func main() {
	sqlHost := os.Getenv("PSQL_HOST")
	sqlUser := os.Getenv("PSQL_USER")
	sqlPassword := os.Getenv("PSQL_PASSWORD")
	sqlDatabase := os.Getenv("PSQL_DATABASE")
	sqlPort := os.Getenv("PSQL_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", sqlHost, sqlUser, sqlPassword, sqlDatabase, sqlPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("could not connect to postgres: ", err)
	}

	err = db.AutoMigrate(&data.PatientDataHeader{}, &data.PatientDataEntry{})
	if err != nil {
		log.Fatalln("could not migrate: ", err)
	}

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Println("invalid method: ", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("error reading body: ", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		patient, err := data.NewPatientFromBytes(bodyBytes)
		if err != nil {
			log.Println("error unmarshalling patient: ", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var patientRecord data.PatientDataHeader
		db.Where("patient_id = ? AND day_number = ?", patient.PatientId, patient.DayNumber).First(&patientRecord)
		if patientRecord.Id == uuid.Nil {
			tx := db.Begin()
			patientRecord = data.PatientDataHeader{
				Id:        uuid.New(),
				PatientId: patient.PatientId,
				DayNumber: patient.DayNumber,
			}
			tx.Create(patientRecord)
			if tx.Error != nil {
				tx.Rollback()
				log.Println("error creating patient record: ", db.Error)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			for _, item := range patient.Items {
				tx.Create(&data.PatientDataEntry{
					Id:                  uuid.New(),
					PatientDataHeaderId: patientRecord.Id,
					Offset:              item.Offset,
					Length:              item.Length,
					Mean:                item.Mean,
				})
				if tx.Error != nil {
					tx.Rollback()
					log.Println("error creating patient record: ", db.Error)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
			tx.Commit()
		} else {
			log.Println("patient record already exists")
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})

	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			log.Println("invalid method: ", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		idParam := r.URL.Query().Get("id")
		dayNumberParam := r.URL.Query().Get("daynumber")
		fmt.Println("id param: ", idParam)
		fmt.Println("day number param: ", dayNumberParam)
		if idParam == "" || dayNumberParam == "" {
			log.Println("missing id or day number param")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var patientRecord data.PatientDataHeader
		db.Model(&data.PatientDataHeader{}).Preload("Items").Where("patient_id = ? AND day_number = ?", idParam, dayNumberParam).First(&patientRecord)
		if patientRecord.Id == uuid.Nil {
			log.Println("patient record not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		marshal, err := json.Marshal(patientRecord)
		if err != nil {
			log.Println("error marshalling patient: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(marshal)
		if err != nil {
			log.Println("error writing response: ", err)
			return
		}
	})

	err = http.ListenAndServe(":5000", nil)
	if err != nil {
		panic(err)
	}
}
