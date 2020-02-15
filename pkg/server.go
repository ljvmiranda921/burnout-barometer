// Copyright 2020 Lester James V. Miranda. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// Server implements the Optserve server to be run inside the cluster.
type Server struct {
	Port   int
	Router *httprouter.Router
	Config *Configuration

	database Database

	// If set to true, then the Request.Process() method will not insert into
	// the database. The resulting Message is just returned.
	DebugOnly bool
}

// Routes contain all handler functions that responds to GET or POST requests.
func (s *Server) Routes() {
	log.Debug("serving routes")
	s.Router.HandlerFunc(http.MethodPost, "/log", s.handleLog())
	s.Router.HandlerFunc(http.MethodGet, "/", s.handleIndex())
}

// Start command starts a server on the specific port.
func (s *Server) Start() error {
	log.Infof("listening to port %d", s.Port)

	// Identify the scheme for the database connection
	db, err := NewDatabase(s.Config.Table)
	if err != nil {
		return err
	}
	s.database = db

	http.ListenAndServe(fmt.Sprintf(":%d", s.Port), s.Router)
	return nil
}

func (s *Server) handleIndex() http.HandlerFunc {
	type response struct {
		Message string `json:"message"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{"path": "/"}).Trace("received request")
		res := response{Message: "PONG"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&res)
	}
}

func (s *Server) handleLog() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{"path": "/log"}).Trace("received request")
		w.Header().Set("Content-Type", "application/json")

		if err := r.ParseForm(); err != nil {
			e := errorMsg{
				Message: fmt.Sprintf("couldn't parse form: %s", err),
				Code:    http.StatusBadRequest,
			}
			e.JSONError(w)
			log.WithFields(log.Fields{"err": e}).Error("http.Request.ParseForm")
			return
		}

		// Check if authentication token is correct
		if err := VerifyWebhook(r.Form, s.Config.Token); err != nil {
			e := errorMsg{
				Message: fmt.Sprintf("token may be missing or invalid: %s", err),
				Code:    http.StatusUnauthorized,
			}
			e.JSONError(w)
			log.WithFields(log.Fields{"err": e.Message}).Error("VerifyWebhook")
			return
		}

		// Check if text exists
		if len(r.Form["text"]) == 0 {
			e := errorMsg{
				Message: "empty text in form",
				Code:    http.StatusBadRequest,
			}
			e.JSONError(w)
			log.Error(e.Message)
			return
		}

		// Create Twitter client
		var client *twitter.Client
		if tc := s.Config; containsEmpty(
			tc.TwitterConsumerKey,
			tc.TwitterConsumerSecret,
			tc.TwitterAccessKey,
			tc.TwitterAccessSecret,
		) {
			client = nil
		} else {
			config := &clientcredentials.Config{
				ClientID:     s.Config.TwitterConsumerKey,
				ClientSecret: s.Config.TwitterConsumerSecret,
				TokenURL:     "https://api.twitter.com/oauth2/token",
			}
			httpClient := config.Client(oauth2.NoContext)
			client = twitter.NewClient(httpClient)
		}

		// Process the request
		req := &Request{
			Text:          r.FormValue("text"),
			UserID:        r.FormValue("user_id"),
			Timestamp:     r.Header.Get("X-Slack-Request-Timestamp"),
			Area:          s.Config.Area,
			DB:            s.database,
			TwitterClient: client,
			DebugOnly:     s.DebugOnly,
		}

		resp, err := req.Process()
		if err != nil {
			e := errorMsg{
				Message: fmt.Sprintf("error in processing request: %s", err),
				Code:    http.StatusBadRequest,
			}
			e.JSONError(w)
			log.WithFields(log.Fields{"err": e.Message}).Error("Request.Process")
			return
		}

		// Send reply back to Slack
		json.NewEncoder(w).Encode(resp)
	}
}

func containsEmpty(ss ...string) bool {
	for _, s := range ss {
		if s == "" {
			return true
		}
	}
	return false
}

// VerifyWebhook checks if the submitted request matches the token provided by Slack
func VerifyWebhook(form url.Values, token string) error {
	t := form.Get("token")
	if len(t) == 0 {
		return fmt.Errorf("empty form token")
	}

	if t != token {
		return fmt.Errorf("invalid request/credentials: %q", t[0])
	}

	return nil
}

type errorMsg struct {
	Message string `json:"message"`
	Code    int    `json:"status_code"`
}

func (e errorMsg) JSONError(w http.ResponseWriter) string {
	w.WriteHeader(e.Code)
	json.NewEncoder(w).Encode(&e)
	return e.Message
}
