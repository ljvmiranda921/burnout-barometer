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
			}
		})
	}
}

func TestParseMessage(t *testing.T) {
	tests := []struct {
		name, arg, want, want1 string
	}{
		{name: "single notes", arg: "4 hello", want: "4", want1: "hello"},
		{name: "multiple notes", arg: "4 hello world", want: "4", want1: "hello world"},
		{name: "no notes", arg: "4", want: "4", want1: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ParseMessage(tt.arg)
			if got != tt.want {
				t.Errorf("Request.parseMessage() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Request.parseMessage() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
