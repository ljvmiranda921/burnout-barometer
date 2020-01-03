// Copyright 2019 Lester James V. Miranda. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

// Server implements the Optserve server to be run inside the cluster.
type Server struct {
	Port   int
	Router *httprouter.Router
}

// Routes contain all handler functions that responds to GET or POST requests.
func (s *Server) Routes() {
	log.Debug("serving routes")
	s.Router.HandlerFunc(http.MethodPost, "/log", s.handleLog())
	s.Router.HandlerFunc(http.MethodGet, "/", s.handleIndex())
}

// Start command starts a server on the specific port. This will always return
// a nil error.
func (s *Server) Start() error {
	log.Infof("listening to port %d", s.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", s.Port), s.Router)
	return nil
}

func (s *Server) handleIndex() http.HandlerFunc {
	type response struct {
		Data string `json:"data"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{"path": "/"}).Trace("received request")
		res := response{Data: "PONG"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&res)
	}
}

func (s *Server) handleLog() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{"path": "/log"}).Trace("received request")

		// Check if authentication token is correct
		if err := verifyWebHook(r.Form); err != nil {
			log.WithFields(log.Fields{"err": err}).Fatal("verifyWebHook")
		}

		// Check if text exists
		if len(r.Form["text"]) == 0 {
			log.Fatal("empty text in form")
		}

		// Process the request
		req := &Request{
			Text:      r.Form["text"][0],
			UserID:    r.Form["user_id"][0],
			Timestamp: r.Header.Get("X-Slack-Request-Timestamp"),
			Area:      config.Area,
			BQTable:   config.Table,
		}

		resp, err := req.Process()
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Fatal("req.Process()")
		}

		// Send reply back to Slack
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
