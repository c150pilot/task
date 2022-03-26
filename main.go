package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Name struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type JokeResponseData struct {
	Type  string `json:"type"`
	Value Joke   `json:"value"`
}

type Joke struct {
	ID         int      `json:"id"`
	Joke       string   `json:"joke"`
	Categories []string `json:"categories"`
}

// Define Client
var client http.Client

func main() {
	// Setup Client & Server
	client = http.Client{
		Timeout: 10 * time.Second,
	}

	server := &http.Server{
		Addr:           ":5000",
		Handler:        http.HandlerFunc(handler),
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("Server Started on Port: 5000")
	log.Fatal(server.ListenAndServe())
}

// Fetch Random Name
func getName() (Name, error) {
	// Declare Return Variable
	var newName Name

	// Fetch First & Last Name from API
	resp, err := client.Get("https://names.mcquay.me/api/v0/")

	// Error Checking
	if err != nil {
		log.Panicln(err)
		return newName, err
	}

	// Defer closing of response body
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Panicln(err)
			return
		}
	}(resp.Body)

	// Read from Response Body
	body, err := ioutil.ReadAll(resp.Body)

	// Error Checking
	if err != nil {
		log.Panicln(err)
		return newName, err
	}

	// Unmarshal JSON string into Name type
	err = json.Unmarshal(body, &newName)

	// Error Handling
	if err != nil {
		log.Panicln(err)
	}

	// Return Name
	return newName, nil
}

// Fetch Random Joke
func getJoke(name Name) (string, error) {
	// Declare Return Variable
	var responseData JokeResponseData

	// Format Get Request
	req := "http://api.icndb.com/jokes/random?firstName=" + name.FirstName + "&lastName=" + name.LastName + "&limitTo=nerdy"

	// Fetch First & Last Name from API
	resp, err := client.Get(req)

	// Error Checking
	if err != nil {
		log.Panicln(err)
		return "", err
	}

	// Defer closing of response body
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Panicln(err)
			return
		}
	}(resp.Body)

	// Read from Response Body
	body, err := ioutil.ReadAll(resp.Body)

	// Error Checking
	if err != nil {
		log.Panicln(err)
		return "", err
	}

	// Unmarshal JSON string into Joke Response Data Type
	err = json.Unmarshal(body, &responseData)

	// Error Checking
	if err != nil {
		log.Panicln(err)
	}

	// Return Joke as string
	return responseData.Value.Joke, nil
}

// Combine Joke and Name
func makeJoke() (string, error) {
	// Get Name for Joke
	name, err := getName()

	// Error Checking
	if err != nil {
		log.Panicln(err)
		return "", err
	}

	// Get Joke
	joke, err := getJoke(name)
	if err != nil {
		log.Panicln(err)
		return "", err
	}

	return joke, nil
}

// Handle Incoming HTTP Requests
func handler(w http.ResponseWriter, r *http.Request) {
	var joke string
	var err error

	// Get Joke
	joke, err = makeJoke()

	// Error Handling
	if err != nil {
		log.Panicln(err)
		return
	}

	fmt.Fprintf(w, joke)
}
