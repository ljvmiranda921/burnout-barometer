// Copyright 2019 Lester James V. Miranda. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

// Package pkg contains types and methods for interacting with the barometer.
package pkg

import (
	"reflect"
	"testing"
)

func TestRequest_Process(t *testing.T) {
	type fields struct {
		Text      string
		UserID    string
		Timestamp string
		Area      string
		DB        Database
		DebugOnly bool
		item      logItem
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Message
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Request{
				Text:      tt.fields.Text,
				UserID:    tt.fields.UserID,
				Timestamp: tt.fields.Timestamp,
				Area:      tt.fields.Area,
				DB:        tt.fields.DB,
				DebugOnly: tt.fields.DebugOnly,
				item:      tt.fields.item,
			}
			got, err := r.Process()
			if (err != nil) != tt.wantErr {
				t.Errorf("Request.Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Request.Process() = %v, want %v", got, tt.want)
			}
		})
	}
}
