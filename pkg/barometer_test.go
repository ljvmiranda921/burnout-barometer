// Copyright 2020 Lester James V. Miranda. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

// Package pkg contains types and methods for interacting with the barometer.
package pkg

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	// Ignore log messages
	log.SetOutput(ioutil.Discard)
}

func TestRequest_Process(t *testing.T) {
	type fields struct {
		Text      string
		UserID    string
		Timestamp string
		Area      string
		DB        Database
		Debug     bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Message
		wantErr bool
	}{
		{
			name:    "happy path",
			fields:  fields{Text: "4 hello world", Debug: true, Timestamp: strconv.FormatInt(time.Now().Unix(), 10), Area: "Asia/Manila"},
			want:    &Message{Text: fmt.Sprintf("%s: 4 (hello world)", ackPrefix)},
			wantErr: false,
		},
		{
			name:    "unknown location",
			fields:  fields{Text: "4 hello world", Debug: true, Timestamp: strconv.FormatInt(time.Now().Unix(), 10), Area: "Europe/Manila"},
			want:    &Message{},
			wantErr: true,
		},
		{
			name:    "unknown timestamp",
			fields:  fields{Text: "4 hello world", Debug: true, Timestamp: "03149a", Area: "Asia/Manila"},
			want:    &Message{},
			wantErr: true,
		},
		{
			name:    "non-int measure",
			fields:  fields{Text: "A hello world", Debug: true, Timestamp: strconv.FormatInt(time.Now().Unix(), 10), Area: "Asia/Manila"},
			want:    &Message{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Request{
				Text:      tt.fields.Text,
				UserID:    tt.fields.UserID,
				Timestamp: tt.fields.Timestamp,
				Area:      tt.fields.Area,
				DB:        tt.fields.DB,
				Debug:     tt.fields.Debug,
			}

			got, err := r.Process()
			if (err != nil) != tt.wantErr {
				t.Errorf("Request.Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) && (got.Text != tt.want.Text) {
				t.Errorf("Request.Process() = %v, want %v", got.Text, tt.want.Text)
			}
		})
	}
}

func TestRequest_parseMessage(t *testing.T) {
	type fields struct {
		Text      string
		UserID    string
		Timestamp string
		Area      string
		DB        Database
		Debug     bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
		want1  string
	}{
		{name: "single notes", fields: fields{Text: "4 hello"}, want: "4", want1: "hello"},
		{name: "multiple notes", fields: fields{Text: "4 hello world"}, want: "4", want1: "hello world"},
		{name: "no notes", fields: fields{Text: "4"}, want: "4", want1: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Request{
				Text:      tt.fields.Text,
				UserID:    tt.fields.UserID,
				Timestamp: tt.fields.Timestamp,
				Area:      tt.fields.Area,
				DB:        tt.fields.DB,
				Debug:     tt.fields.Debug,
			}
			got, got1 := r.message()
			if got != tt.want {
				t.Errorf("Request.parseMessage() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Request.parseMessage() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
