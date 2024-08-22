package resources

import (
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Submission struct {
	answer  string
	when    time.Time
	correct bool
	message string
}

func NewSubmission(data *http.Response, answer string) (*Submission, error) {
	defer data.Body.Close()
	doc, err := goquery.NewDocumentFromReader(data.Body)
	if err != nil {
		return nil, err
	}

	message := doc.Find("article").Text()

	newSub := &Submission{
		when:    time.Now(),
		answer:  answer,
		message: message,
	}

	if strings.Contains(message, "right answer") || strings.Contains(message, "Congratulations!") {
		newSub.correct = true
	} else {
		newSub.correct = false
	}

	return newSub, nil
}
