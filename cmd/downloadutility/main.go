package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		str := fmt.Sprintf("Missing patient id or day number. Usage: %s <patient id> <day number>", args[0])
		log.Fatalln(str)
	}

	patientId := args[1]
	dayNumber := args[2]

	uri := fmt.Sprintf("http://localhost:5000/download?id=%s&daynumber=%s", patientId, dayNumber)
	request, err := http.Get(uri)
	if err != nil {
		log.Fatalln("could not create request: ", err)
	}
	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		log.Fatalln("could not read body: ", err)
	}

	fmt.Println(string(bodyBytes))
}
