// Copyright 2020 Lester James V. Miranda. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

package pkg

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	// Ignore log messages
	log.SetOutput(ioutil.Discard)
}

func TestUpdateLog(t *testing.T) {
	type args struct {
		userID, text string
		timestamp    time.Time
		db           DBInserter
		debug        bool
	}
	tests := []struct {
		name    string
		args    args
		want    *Message
		wantErr bool
	}{
		{
			name:    "happy path",
			args:    args{text: "4 hello world", debug: true, timestamp: time.Now()},
			want:    &Message{Text: fmt.Sprintf("%s: 4 (hello world)", ackPrefix)},
			wantErr: false,
		},
		{
			name:    "non-int measure",
			args:    args{text: "A hello world", debug: true, timestamp: time.Now()},
			want:    &Message{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UpdateLog(tt.args.userID, tt.args.text, tt.args.timestamp, tt.args.db, nil, tt.args.debug)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateLog() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) && (got.Text != tt.want.Text) {
				t.Errorf("UpdateLog() = %v, want %v", got.Text, tt.want.Text)
				return
			}
		})
	}
}

func TestParseMessage(t *testing.T) {
	tests := []struct {
		name, arg, want1 string
		want             int
		wantErr          bool
	}{
		{name: "single notes", arg: "4 hello", want: 4, want1: "hello", wantErr: false},
		{name: "multiple notes", arg: "4 hello world", want: 4, want1: "hello world", wantErr: false},
		{name: "no notes", arg: "4", want: 4, want1: "", wantErr: false},
		{name: "cannot convert measure", arg: "X hello world", wantErr: true},
		{name: "float measure", arg: "2.0 hello world", wantErr: true},
		{name: "log measure outside range 1", arg: "100 hello world", wantErr: true},
		{name: "log measure outside range 2", arg: "-100 hello world", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ParseMessage(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) && *got != tt.want {
				t.Errorf("ParseMessage() got = %d, want %d", *got, tt.want)
				return
			}
			if (got1 != nil) && *got1 != tt.want1 {
				t.Errorf("ParseMessage() got1 = %s, want %s", *got1, tt.want1)
				return
			}
		})
	}
}

func ExampleUpdateLog() {
	// Prepare inputs for updating the log
	userID := "W012A3CDE"
	text := "4 Had dinner with friends today!"
	message, err := UpdateLog(userID, text, time.Now(), nil, nil, true) // Run in debug-mode
	if err != nil {
		log.Fatalf("cannot update log, err: %v", err)
	}
	fmt.Println(message.Text)
	// Output: Gotcha, I logged your mood: 4 (Had dinner with friends today!)
}

func ExampleParseMessage() {
	message := "4 Had awesome dinner!"
	measure, notes, err := ParseMessage(message)
	if err != nil {
		log.Fatalf("cannot parse message, err: %v", err)
	}
	fmt.Printf("Your message: %s (%d)", *notes, *measure)
	// Output: Your message: Had awesome dinner! (4)
}
