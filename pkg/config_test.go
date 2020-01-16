// Copyright 2019 Lester James V. Miranda. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

package pkg

import (
	"os"
	"testing"
)

func TestConfiguration_WriteConfiguration(t *testing.T) {
	type fields struct {
		ProjectID string
		Table     string
		Token     string
		Area      string
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
			fields:  fields{ProjectID: "test-project", Table: "test-table", Token: "test-token", Area: "test-area"},
			args:    args{outputPath: "test_file.json"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Configuration{
				ProjectID: tt.fields.ProjectID,
				Table:     tt.fields.Table,
				Token:     tt.fields.Token,
				Area:      tt.fields.Area,
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
			want:    &Configuration{ProjectID: "test-project", Table: "bq://test-table", Token: "ZK[VPIHE9E2CIMAz0QUE", Area: "Asia/Manila"},
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
			if tt.want.ProjectID != got.ProjectID {
				t.Errorf("ReadConfiguration() = %v, want %v", got.ProjectID, tt.want.ProjectID)
			}
		})
	}
}
