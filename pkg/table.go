package pkg

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"cloud.google.com/go/bigquery"
)

// Database is a generic interface for storing and accessing barometer logs.
type Database interface {
	Insert(item Log) error
}

// BigQuery provides a connection to BigQuery. It implements the Table
// interface.
type BigQuery struct {
	URL string

	options                 url.URL
	project, dataset, table string
}

// Insert adds an Item entry into the specified BigQuery table
func (t *BigQuery) Insert(item Log) error {
	ctx := context.Background()
	project, dataset, table := t.splitBQPath(t.options.Host)
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

func (t *BigQuery) splitBQPath(p string) (string, string, string) {
	s := strings.Split(p, ".")
	return s[0], s[1], s[2]
}
