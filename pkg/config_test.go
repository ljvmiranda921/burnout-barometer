// Copyright 2020 Lester James V. Miranda. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

package pkg

import (
	"fmt"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestConfiguration_WriteConfiguration(t *testing.T) {
	type fields struct {
		Table string
		Token string
		Area  string
	}
	type args struct {
		outputPath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "file exists after creation",
			fields:  fields{Table: "test-table", Token: "test-token", Area: "test-area"},
			args:    args{outputPath: "test_file.json"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Configuration{
				Table: tt.fields.Table,
				Token: tt.fields.Token,
				Area:  tt.fields.Area,
			}
			if err := cfg.WriteConfiguration(tt.args.outputPath); (err != nil) != tt.wantErr {
				t.Errorf("Configuration.WriteConfiguration() error = %v, wantErr %v", err, tt.wantErr)
			}
			// Remove the generated file
			if _, err := os.Stat(tt.args.outputPath); err == nil {
				t.Logf("removing generated file: %s", tt.args.outputPath)
				os.Remove(tt.args.outputPath)

			}
		})
	}
}

func TestReadConfiguration(t *testing.T) {
	type args struct {
		cfgPath string
	}
	tests := []struct {
		name    string
		args    args
		want    *Configuration
		wantErr bool
	}{
		{
			name:    "happy path read config",
			args:    args{cfgPath: "testdata/test_happy_path_read_config.json"},
			want:    &Configuration{Table: "bq://test-table", Token: "ZK[VPIHE9E2CIMAz0QUE", Area: "Asia/Manila"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := ReadConfiguration(tt.args.cfgPath)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadConfiguration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func ExampleReadConfiguration() {
	// Read config from a file
	config, err := ReadConfiguration("path/to/config.json")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v", config)
}

func ExampleConfiguration_WriteConfiguration() {
	// Create a sample configuration
	config := &Configuration{
		Table: "bq://my-project.my-dataset.my-table",
		Token: "M4KY3LOVPIhE9E2zIMAz0QUE",
		Area:  "Asia/Manila",
	}

	err := config.WriteConfiguration("path/to/config.json")
	if err != nil {
		log.Fatal(err)
	}

}
