package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"golang.org/x/time/rate"
)

const USER_AGENT = "dalton.dog/aocgo"
const BASE_URL = "https://adventofcode.com"
const YEAR_URL = BASE_URL + "/%v"
const DAY_URL = YEAR_URL + "/day/%v"

const REQS_PER_SEC = 20

var MasterClient httpClient

type httpClient struct {
	client       http.Client
	sessionToken string // Eventually make this []string in case we want to run for multiple users?
	rateLimiter  *rate.Limiter
}

// Do is a wrapper around a normal client call in
// order to use our rate limiter.
func (c *httpClient) Do(req *http.Request) (*http.Response, error) {
	err := c.rateLimiter.Wait(context.Background())
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// InitClient creates the API client with a given user session token
func InitClient(userSessionToken string) {
	limiter := rate.NewLimiter(rate.Every(time.Second), REQS_PER_SEC)
	client := httpClient{
		client:       http.Client{},
		sessionToken: userSessionToken,
		rateLimiter:  limiter,
	}
	MasterClient = client
}

func NewGetReq(url string, sessionToken string) (*http.Response, error) {
	log.Debug("Making GET request.", "URL", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Error creating GET request!", "error", err)
	}

	if sessionToken == "" {
		sessionToken = MasterClient.sessionToken
	}

	req.Header.Add("User-Agent", USER_AGENT)
	req.Header.Add("Cookie", fmt.Sprintf("session=%s", strings.TrimSpace(sessionToken)))

	return MasterClient.Do(req)
}

// NOTE: To submit answers:
//	Method:  POST
//	URL:	 https://adventofcode.com/yyyy/day/d/answer
//	Headers:
//		Cookie: session=<session token>
//		Content-Type: application/x-www-form-urlencoded
//	Form Data (Body):
//		`level` : 1 if Part A, 2 if Part B
//		`answer` : Answer to submit

func SubmitAnswer(year int, day int, part int, userSession string, answer string) (*http.Response, error) {
	URL := PuzzleAnswerURL(year, day)
	log.Debugf("Attempting to submit answer for Day %v (%v) [Part %v] to URL %v", day, year, part, URL)
	log.Debugf("Answer: %v -- User: %v", answer, userSession)

	formData := url.Values{}
	formData.Set("level", strconv.Itoa(part))
	formData.Set("answer", answer)

	encodedForm := formData.Encode()

	req, err := http.NewRequest("POST", URL, strings.NewReader(encodedForm))
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", USER_AGENT)
	req.Header.Add("Cookie", fmt.Sprintf("session=%v", userSession))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return MasterClient.client.Do(req)
}

// URL Helper Methods

func PuzzlePageURL(year int, day int) string {
	return fmt.Sprintf(DAY_URL, year, day)
}

func PuzzleInputURL(year int, day int) string {
	return fmt.Sprintf(DAY_URL, year, day) + "/input"
}

func PuzzleAnswerURL(year int, day int) string {
	return fmt.Sprintf(DAY_URL, year, day) + "/answer"
}
