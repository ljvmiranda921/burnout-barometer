package pkg

import (
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestRequest_ParseMessage(t *testing.T) {
	type test struct {
		input string
		want  []string
	}

	tests := []test{
		{input: "4 hello", want: []string{"4", "hello"}},
		{input: "4 hello world", want: []string{"4", "hello world"}},
		{input: "4", want: []string{"4", ""}},
	}

	for _, tc := range tests {
		request := &Request{Text: tc.input}
		measure, notes, _ := request.ParseMessage()
		got := []string{measure, notes}
		if !reflect.DeepEqual(tc.want, got) {
			t.Errorf("expected: %v, got: %v", tc.want, got)
		}
	}
}

func TestRequest_GetTimestamp(t *testing.T) {
	timeNow := strconv.FormatInt(time.Now().Unix(), 10)
	t.Logf("time now in unix: %s", timeNow)
	request := &Request{Timestamp: timeNow, Area: "Asia/Manila"}
	ts, err := request.GetTimestamp()
	if err != nil {
		t.Errorf("error getting time: %v", err)
	}
	t.Logf("time now in human-readable format: %s", ts)
}
