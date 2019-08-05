package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"testing"
)

func sendPost(bytesRepresentation []byte) {
	resp, err := http.Post("http://localhost:8080", "application/json", bytes.NewBuffer(bytesRepresentation))

	if err != nil {
		log.Fatalln(err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%s", b)
}

func TestPost(t *testing.T) {
	for i := 1; i < 10; i++ {
		id := strconv.Itoa(i)
		jsonBody := map[string]string{
			"Name":  "User" + id,
			"Email": "user" + id + "@semrush.com",
		}

		bytesRepresentation, err := json.Marshal(jsonBody)
		if err != nil {
			log.Fatalln(err)
		}

		go sendPost(bytesRepresentation)
	}
}
