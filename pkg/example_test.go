package pkg

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

func ExampleParseMessage() {
	message := "4 Had awesome dinner!"
	measure, notes := ParseMessage(message)
	fmt.Printf("Your message: %s (%s)", notes, measure)
	// Output: Your message: Had awesome dinner! (4)
}

func ExampleFetchTimestamp() {
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	area := "Asia/Manila" // IANA-compliant timezone name
	timestamp, err := FetchTimestamp(ts, area)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("time is: %v", timestamp)
}

func ExampleContainsEmpty() {
	containsEmpty := []string{"A", "B", ""}
	if ContainsEmpty(containsEmpty...) {
		fmt.Println("contains an empty string")
	}
	// Output:  contains an empty string
}

func ExampleVerifyWebhook() {
	// Example token obtained from Slack
	token := "M4KY3LOVPIhE9E2zIMAz0QUE"

	// Example query sent by the slash command
	q := "token=M4KY3LOVPIhE9E2zIMAz0QUE&text=4 hello&user_id=UA1DXYCL2"
	v, err := url.ParseQuery(q)
	if err != nil {
		panic(err)
	}

	if err := VerifyWebhook(v, token); err != nil {
		// Webhook didn't match, throw an error
		log.Fatal(err)
	}
}
