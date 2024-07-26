package internal

import (
	"fmt"
	"io"
	"net/http"

	"dalton.dog/aocutil/internal/models"
)

const USER_AGENT = "dalton.dog/aocutil/0.0"
const BASE_URL = "https://adventofcode.com"
const YEAR_URL = BASE_URL + "/%v"
const DAY_URL = YEAR_URL + "/day/%v"

// WARN: Be sure to implement rate limiting from the start. Try to make access as efficient as possible
//		https://github.com/wimglenn/advent-of-code-data/issues/59

// TODO: GetData()

// TODO: SubmitGuess()

var MasterClient *httpClient

type httpClient struct {
	client       *http.Client
	sessionToken string
}

func InitClient(userSessionToken string) {
	client := &httpClient{
		client:       &http.Client{},
		sessionToken: userSessionToken,
	}
	MasterClient = client
}

func newGetReq(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
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
//	Form Data (Body):
//		`level` : 1 if Part A, 2 if Part B
//		`answer` : Answer to submit

// My `answer` tests in Postman weren't working. Not sure what I was doing wrong. Maybe it'll work here
// So... Postman at home was working fine? I am confuse, but oh well lol
func NewPostReq() {}

func GetGenericPuzzleData(day int, year int) {
	URL := fmt.Sprintf(DAY_URL, year, day)
	resp, err := newGetReq(URL)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}

func GetUserPuzzleInput(year int, day int, userSession string) []byte {
	resp, err := newGetReq(PuzzleInputURL(year, day))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	inputData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return inputData
}

func GetData(user *models.User, day int, year int) string {
	return ""
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
