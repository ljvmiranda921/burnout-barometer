// Copyright 2019 Lester James V. Miranda. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

package pkg

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"cloud.google.com/go/bigquery"
	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
)

// Database is a generic interface for storing and accessing barometer logs.
type Database interface {
	GetURL() *url.URL
	Insert(item Log) error
}

// NewDatabase creates a Database based on the detected scheme of the URL.
func NewDatabase(dburl string) (Database, error) {
	u, err := url.Parse(dburl)
	if err != nil {
		return nil, err
	}

	var db Database

	switch u.Scheme {
	case "bigquery", "bq":
		log.WithFields(log.Fields{"scheme": u.Scheme}).Info("detected scheme")
		db = &bigQuery{URL: dburl, Config: u}
	default:
		msg := fmt.Sprintf("unknown database scheme: %s", u.Scheme)
		log.Fatal(msg)
	}

	return db, nil
}

// BigQuery

type bigQuery struct {
	URL    string
	Config *url.URL

	project, dataset, table string
}

func (t *bigQuery) GetURL() *url.URL {
	return t.Config
}

func (t *bigQuery) Insert(item Log) error {
	ctx := context.Background()
	project, dataset, table := t.splitBQPath(t.Config.Host)
	client, err := bigquery.NewClient(ctx, project)
	if err != nil {
		return fmt.Errorf("error in bigquery.NewClient: %v", err)
	}

	inserter := client.Dataset(dataset).Table(table).Inserter()
	items := []*Log{&item}

	if err := inserter.Put(ctx, items); err != nil {
		return err
	}
	return nil
}

func (t *bigQuery) splitBQPath(p string) (string, string, string) {
	s := strings.Split(p, ".")
	return s[0], s[1], s[2]
}

// Postgres

type postgres struct {
	URL    string
	Config *url.URL
}

func (t *postgres) GetURL() *url.URL {
	return t.Config
}

func (t *postgres) Insert(item Log) error {
	opts, err := pg.ParseURL(t.URL)
	if err != nil {
		return fmt.Errorf("error in pg.ParseURL: %v", err)
	}

	db := pg.Connect(opts)
	defer db.Close()

	if err := db.Insert(&item); err != nil {
		return fmt.Errorf("error in db.Insert: %v", err)
	}

	return nil
}
