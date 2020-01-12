package pkg

import (
	"reflect"
	"testing"
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
