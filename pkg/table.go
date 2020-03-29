// Copyright 2020 Lester James V. Miranda. All rights reserved.
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

// DBInserter is an interface for storing barometer logs.
type DBInserter interface {
	InsertDB(item LogItem) error // Insert a log into the Database
}

// NewDBInserter creates a DBInserter based on the detected scheme of the URL.
func NewDBInserter(dburl string) (DBInserter, error) {
	u, err := url.Parse(dburl)
	if err != nil {
		return nil, err
	}

	var db DBInserter

	switch u.Scheme {
	case "bigquery", "bq":
		log.WithFields(log.Fields{"scheme": u.Scheme}).Info("detected scheme")
		db = &bigQuery{URL: dburl, Config: u}
	case "postgres":
		log.WithFields(log.Fields{"scheme": u.Scheme}).Info("detected scheme")
		db = &postgres{URL: dburl, Config: u}
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
}

func (t *bigQuery) InsertDB(item LogItem) error {
	ctx := context.Background()
	project, dataset, table := t.splitBQPath(t.Config.Host)
	client, err := bigquery.NewClient(ctx, project)
	if err != nil {
		return fmt.Errorf("error in bigquery.NewClient: %v", err)
	}

	inserter := client.Dataset(dataset).Table(table).Inserter()
	items := []*LogItem{&item}

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

func (t *postgres) InsertDB(item LogItem) error {
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
