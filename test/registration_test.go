package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"testing"
	"time"
)

const (
	url              = "http://localhost:5001/api/v1/users"
	totalRequests    = 15
	concurrencyLimit = 10 // Adjust this to control concurrency
)

func generateRandomEmail() string {
	return fmt.Sprintf("%d-user-%d@example-%d.com", rand.Intn(1000000), rand.Intn(1000000), rand.Intn(1000000))
}

type RequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RequestResult struct {
	Email  string
	Status int
	Error  error
}

func sendRequest(client *http.Client, email string, wg *sync.WaitGroup, semaphore chan struct{}, results chan<- RequestResult) {
	defer wg.Done()
	defer func() { <-semaphore }()

	body := RequestBody{
		Email:    email,
		Password: "this-sample-password",
	}

	jsonData, _ := json.Marshal(body)
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		results <- RequestResult{Email: email, Status: 0, Error: err}
		return
	}
	defer resp.Body.Close()

	results <- RequestResult{Email: email, Status: resp.StatusCode, Error: nil}
}

func checkServerRunning() bool {
	client := &http.Client{
		Timeout: 1 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return true
}

func TestSendRequest(t *testing.T) {
	if !checkServerRunning() {
		fmt.Printf("\n⚠️  API Server is not running at %s - Skipping integration tests\n", url)
		t.Skip("Skipping integration test because server is not running")
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrencyLimit)
	client := &http.Client{}
	results := make(chan RequestResult, totalRequests)

	for i := 0; i < totalRequests; i++ {
		email := generateRandomEmail()
		wg.Add(1)
		semaphore <- struct{}{} // Limit concurrency
		go sendRequest(client, email, &wg, semaphore, results)
	}

	wg.Wait()
	close(results)

	totalFailed := 0
	totalSucceed := 0
	for result := range results {
		if result.Error != nil {
			totalFailed++
			t.Errorf("Request failed for %s: %v", result.Email, result.Error)
		} else if result.Status < 200 || result.Status >= 300 {
			totalFailed++
			t.Errorf("Unexpected status code %d for %s", result.Status, result.Email)
		} else {
			totalSucceed++
		}
	}

	fmt.Println("Total failed requests:", totalFailed)
	fmt.Println("Total succeeded requests:", totalSucceed)
}
