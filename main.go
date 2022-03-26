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

type NameResult struct {
	Value Name
	Error error
}

type JokeResult struct {
	Value string
	Error error
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
func getName(result chan NameResult) {
	var newName Name
	var resultData NameResult // To Return

	// Setup to Handle Too Many Requests Error from Random Name API
	isError := true
	var resp *http.Response
	var err error

	for isError {
		// Fetch First & Last Name from API
		resp, err = client.Get("https://names.mcquay.me/api/v0/")

		// Error Checking
		if err != nil {
			log.Panicln(err)
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
		}

		// Unmarshal JSON string into Name type
		err = json.Unmarshal(body, &newName)

		// Error Handling
		if err != nil {
			fmt.Println("Too Many Requests Error from Name API Server - Retrying")
		} else {
			isError = false
		}

	}

	// Add Name to Channel
	resultData = NameResult{newName, err}
	result <- resultData
}

// Fetch Random Joke
func getJoke(name Name, result chan JokeResult) {
	var responseData JokeResponseData
	var resultData JokeResult // To Return

	// Format Get Request
	req := "http://api.icndb.com/jokes/random?firstName=" + name.FirstName + "&lastName=" + name.LastName + "&limitTo=nerdy"

	// Fetch First & Last Name from API
	resp, err := client.Get(req)

	// Error Checking
	if err != nil {
		resultData = JokeResult{responseData.Value.Joke, nil}
		result <- resultData
		log.Panicln(err)
	}

	// Defer closing of response body
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			resultData = JokeResult{responseData.Value.Joke, nil}
			result <- resultData
			log.Panicln(err)
		}
	}(resp.Body)

	// Read from Response Body
	body, err := ioutil.ReadAll(resp.Body)

	// Error Checking
	if err != nil {
		resultData = JokeResult{responseData.Value.Joke, nil}
		result <- resultData
		log.Panicln(err)
	}

	// Unmarshal JSON string into Joke Response Data Type
	err = json.Unmarshal(body, &responseData)

	// Error Checking
	if err != nil {
		resultData = JokeResult{responseData.Value.Joke, nil}
		result <- resultData
		log.Panicln(err)
	}

	// Add Joke to Channel
	resultData = JokeResult{responseData.Value.Joke, nil}
	result <- resultData
}

// Combine Joke and Name
func makeJoke(result chan JokeResult) {
	// Setup Channels to help with concurrency
	nameResult := make(chan NameResult)
	jokeResult := make(chan JokeResult)

	// Get Name for Joke
	go getName(nameResult)

	// Read from nameResult Channel
	nResult := <-nameResult
	err := nResult.Error
	name := nResult.Value

	// Error Checking
	if err != nil {
		resultData := JokeResult{"", err}
		result <- resultData
		log.Panicln(err)
		return
	}

	// Get Joke
	go getJoke(name, jokeResult)

	// Read from JokeResult Channel
	jResult := <-jokeResult
	err = jResult.Error
	joke := jResult.Value

	if err != nil {
		resultData := JokeResult{"", err}
		result <- resultData
		log.Panicln(err)
		return
	}
	resultData := JokeResult{joke, nil}
	result <- resultData
}

// Handle Incoming HTTP Requests
func handler(w http.ResponseWriter, r *http.Request) {
	// Get Joke
	result := make(chan JokeResult)
	go makeJoke(result)
	data := <-result

	// Error Handling
	if data.Error != nil {
		log.Panicln(data.Error)
		return
	}

	fmt.Fprintf(w, data.Value)
}
