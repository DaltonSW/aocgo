package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
)

const USER_AGENT = "dalton.dog/aocutil/0.0"
const BASE_URL = "https://adventofcode.com"
const YEAR_URL = BASE_URL + "/%v"
const DAY_URL = YEAR_URL + "/day/%v"

// WARN: Be sure to implement rate limiting from the start. Try to make access as efficient as possible
//		https://github.com/wimglenn/advent-of-code-data/issues/59

// TODO: GetData()

// TODO: SubmitGuess()

var MasterClient httpClient

type httpClient struct {
	client       http.Client
	sessionToken string // Eventually make this []string in case we want to run for multiple users?
}

func InitClient(userSessionToken string) {
	log.Debug("Initiating API client.", "sessionToken", userSessionToken)
	client := httpClient{
		client:       http.Client{},
		sessionToken: userSessionToken,
	}
	MasterClient = client
}

func NewGetReq(url string) (*http.Response, error) {
	log.Debug("Making GET request.", "URL", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Error creating GET request!", "error", err)
	}

	// We don't NEED to send a User-Agent, but it feels respectful in case we need to get yelled at
	req.Header.Add("User-Agent", USER_AGENT)
	req.Header.Add("Cookie", fmt.Sprintf("session=%v", MasterClient.sessionToken))

	return MasterClient.client.Do(req)
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

// My `answer` tests in Postman weren't working. Not sure what I was doing wrong. Maybe it'll work here
// So... Postman at home was working fine? I am confuse, but oh well lol
func SubmitAnswer(year int, day int, part int, userSession string, answer string) error {
	URL := PuzzleAnswerURL(year, day)
	log.Infof("Attempting to submit answer for Day %v (%v) [Part %v] to URL %v", day, year, part, URL)
	log.Infof("Answer: %v -- User: %v", answer, userSession)

	formData := url.Values{}
	formData.Set("level", strconv.Itoa(part))
	formData.Set("answer", answer)

	encodedForm := formData.Encode()

	req, err := http.NewRequest("POST", URL, strings.NewReader(encodedForm))
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", USER_AGENT)
	req.Header.Add("Cookie", fmt.Sprintf("session=%v", userSession))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := MasterClient.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Infof("Response: %v", string(data))
	return nil
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
