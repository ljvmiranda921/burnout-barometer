// Copyright 2019 Lester James V. Miranda. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

package pkg

import (
	"testing"
)

func TestRequest_ParseMessage(t *testing.T) {
	type fields struct {
		Text      string
		UserID    string
		Timestamp string
		Area      string
		DB        Database
		Item      Log
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		want1   string
		wantErr bool
	}{
		{name: "single notes", fields: fields{Text: "4 hello"}, want: "4", want1: "hello", wantErr: false},
		{name: "multiple notes", fields: fields{Text: "4 hello world"}, want: "4", want1: "hello world", wantErr: false},
		{name: "no notes", fields: fields{Text: "4"}, want: "4", want1: "", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Request{
				Text:      tt.fields.Text,
				UserID:    tt.fields.UserID,
				Timestamp: tt.fields.Timestamp,
				Area:      tt.fields.Area,
				DB:        tt.fields.DB,
				Item:      tt.fields.Item,
			}
			got, got1, err := r.ParseMessage()
			if (err != nil) != tt.wantErr {
				t.Errorf("Request.ParseMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Request.ParseMessage() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Request.ParseMessage() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
