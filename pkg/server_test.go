// Copyright 2020 Lester James V. Miranda. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

package pkg

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

func TestServer_handleIndex(t *testing.T) {
	type fields struct {
		Port     int
		Router   *httprouter.Router
		Config   *Configuration
		database Database
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{name: "happy path", fields: fields{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Port:     tt.fields.Port,
				Router:   tt.fields.Router,
				Config:   tt.fields.Config,
				database: tt.fields.database,
			}

			srv := httptest.NewServer(s.handleIndex())
			defer srv.Close()

			res, err := http.Get(fmt.Sprintf("%s/", srv.URL))
			if err != nil {
				t.Fatalf("could not send GET request: %v", err)
			}
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				t.Errorf("expected status OK; got %v", res.Status)
			}

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}
			t.Logf("received response: %s", string(b))
		})
	}
}

func TestServer_handleLog(t *testing.T) {
	type fields struct {
		Port      int
		Router    *httprouter.Router
		Config    *Configuration
		DebugOnly bool
		database  Database
	}
	type data struct {
		text, userID, token string
	}
	tests := []struct {
		name   string
		data   data
		fields fields
	}{
		{
			name: "happy path",
			data: data{text: "4 hello world", userID: "testUser", token: "testToken"},
			fields: fields{
				Port:      8080,
				DebugOnly: true,
				Config:    &Configuration{Token: "testToken", Area: "Asia/Manila"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Port:      tt.fields.Port,
				Router:    tt.fields.Router,
				Config:    tt.fields.Config,
				DebugOnly: tt.fields.DebugOnly,
				database:  tt.fields.database,
			}

			srv := httptest.NewServer(s.handleLog())
			defer srv.Close()

			client := &http.Client{} // client for sending requests

			// Prepare request
			data := url.Values{}
			data.Add("text", tt.data.text)
			data.Add("user_id", tt.data.userID)
			data.Add("token", tt.data.token)
			payload := strings.NewReader(data.Encode())
			req, err := http.NewRequest("POST", fmt.Sprintf("%s/log", srv.URL), payload)
			if err != nil {
				t.Fatalf("cannot create request: %v", err)
			}
			t.Logf(data.Encode())

			// Add headers
			req.Header.Add("X-Slack-Request-Timestamp", "1579324284")
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			// Perform request
			res, err := client.Do(req)
			if err != nil {
				t.Fatalf("could not send POST request: %v", err)
			}
			defer res.Body.Close()

			// Test proper
			if res.StatusCode != http.StatusOK {
				t.Errorf("expected status OK; got %v", res.Status)
			}

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}
			t.Logf("received response: %s", string(b))
		})
	}
}

func TestVerifyWebhook(t *testing.T) {
	type args struct {
		query string
		token string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "test equal values", args: args{query: "token=sampleToken", token: "sampleToken"}, wantErr: false},
		{name: "test unequal values", args: args{query: "token=notSampleToken", token: "sampleToken"}, wantErr: true},
		{name: "test empty form", args: args{query: "", token: "sampleToken"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			v, err := url.ParseQuery(tt.args.query)
			if err != nil {
				t.Fatalf("cannot parse query: %s", tt.args.query)
			}
			if err := VerifyWebhook(v, tt.args.token); (err != nil) != tt.wantErr {
				t.Errorf("VerifyWebhook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
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
