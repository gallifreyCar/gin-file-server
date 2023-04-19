package utils

import (
	"bytes"
	"os"
	"testing"
)

func TestCreateUploadBody(t *testing.T) {
	files, clean, err := CreateTempFiles(1, 100, "test.txt")
	if err != nil {
		t.Errorf(err.Error())
	}
	defer clean()
	type args struct {
		field string
		files []*os.File
	}
	tests := []struct {
		name    string
		args    args
		want    *bytes.Buffer
		wantErr bool
	}{
		{
			name: "Test single file upload",
			args: args{
				field: "file",
				files: files,
			},
			want:    bytes.NewBufferString("Content-Disposition: form-data; name=\"file\";"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateUploadBody(tt.args.field, tt.args.files)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUploadBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !bytes.Contains(got.Bytes(), tt.want.Bytes()) {
				t.Errorf("CreateUploadBody() got = %v, want %v", got, tt.want)
			}
		})
	}
}
