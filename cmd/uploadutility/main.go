// This program is used to process and send the events to the web service.
// It would emulate an application running on an IoT device.

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"take-home/pkg/data"
)

func readPatientFile(filename string) *data.PatientDto {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln("could not open file: ", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	var patient *data.PatientDto

	patientId := 0
	patientDay := 0
	itemCount := 0
	itemSum := 0
	itemOffset := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineInt, err := strconv.Atoi(line)
		if err != nil {
			log.Fatalln("could not convert line to int: ", err)
		}
		if lineNumber == 0 {
			patientId = lineInt
		} else if lineNumber == 1 {
			patientDay = lineInt
			patient = data.NewPatient(patientId, patientDay)
		} else {
			if lineInt == 0 {
				if itemCount > 0 {
					patient.AddItem(itemOffset, itemCount, itemSum/itemCount)
				}
				itemOffset = lineNumber - 1
				itemSum = 0
				itemCount = 0
			} else {
				itemCount++
				itemSum += lineInt
			}
		}
		lineNumber++
	}
	patient.AddItem(itemOffset, itemCount, itemSum/itemCount)
	return patient
}

func main() {
	args := os.Args
	if len(args) < 2 {
		str := fmt.Sprintf("Missing file name. Usage: %s <file name>", args[0])
		log.Fatalln(str)
	}

	filename := args[1]
	patient := readPatientFile(filename)

	marshal, err := json.Marshal(patient)
	if err != nil {
		log.Fatalln("could not marshal patient: ", err)
	}

	post, err := http.Post("http://localhost:5000/upload", "application/json", bytes.NewReader(marshal))
	if err != nil || post.StatusCode != 201 {
		if post != nil {
			log.Printf("status: %s\n", post.Status)
		}
		log.Fatalln("error uploading data")
	}
}
