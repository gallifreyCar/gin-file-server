package utils

import (
	"reflect"
	"testing"
)

func TestCreateTempFiles(t *testing.T) {
	type args struct {
		fileNum  int
		fileSize int64
		suffix   string
	}
	tests := []struct {
		name         string
		args         args
		wantFilesNum int
		wantErr      bool
	}{
		{
			name: "test create a 10kb file",
			args: args{
				fileNum:  1,
				fileSize: 1024 * 1024 * 10,
				suffix:   "example*.txt",
			},
			wantErr:      false,
			wantFilesNum: 1,
		},
		{
			name: "test create 3 5kb files",
			args: args{
				fileNum:  3,
				fileSize: 1024 * 1024 * 5,
				suffix:   "example*.txt",
			},
			wantErr:      false,
			wantFilesNum: 3,
		},
		{
			name: "test invalid fileNum",
			args: args{
				fileNum:  -3,
				fileSize: 1024 * 1024 * 5,
				suffix:   "example*.txt",
			},
			wantErr:      true,
			wantFilesNum: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFiles, gotClose, err := CreateTempFiles(tt.args.fileNum, tt.args.fileSize, tt.args.suffix)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTempFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(len(gotFiles), tt.wantFilesNum) {
				t.Errorf("CreateTempFiles() gotFiles = %v, wantFilesNum %v", len(gotFiles), tt.wantFilesNum)
			}
			if err == nil {
				defer gotClose()
			}

		})
	}
}
