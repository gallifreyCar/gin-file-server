package respository

import (
	"github.com/gallifreyCar/gin-file-server/model"
	"gorm.io/gorm"
	"reflect"
	"testing"
	"time"
)

func TestGetDataBase(t *testing.T) {
	db, err := getDataBase()
	if err != nil {
		t.Error(err)
	}
	// get sqlDB
	sqlDB, err := db.DB()
	if err != nil {
		t.Error(err)
	}

	// test connection
	if err := sqlDB.Ping(); err != nil {
		t.Error(err)
	}
}

func TestInsertFileLog(t *testing.T) {
	db, err := getDataBase()
	if err != nil {
		t.Error(err)
	}
	id, rowsAffected, err := insertFileLog("testSavePath", "testFileName", "tstUserAgent", "testFileType", db)
	t.Logf("ID: %v,RowsAffected: %v", id, rowsAffected)
	if err != nil {
		t.Error(err)
	}

}

func TestSelectFileLog(t *testing.T) {

	layout := "2006-01-02 15:04:05 -0700 MST"
	loc, _ := time.LoadLocation("Local")
	createAt, _ := time.ParseInLocation(layout, "2023-04-07 16:43:28 +0800 CST", loc)
	updatedAt, _ := time.ParseInLocation(layout, "2023-04-07 16:43:28 +0800 CST", loc)

	db, _ := getDataBase()

	type args struct {
		fileName string
		db       *gorm.DB
	}
	tests := []struct {
		name    string
		args    args
		want    model.UploadFileLog
		wantErr bool
	}{
		{
			name: "Valid file name",
			args: args{
				fileName: "testFileName",
				db:       db,
			},
			want: model.UploadFileLog{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: createAt,
					UpdatedAt: updatedAt,
					DeletedAt: gorm.DeletedAt{},
				},
				SavePath:  "testSavePath",
				FileName:  "testFileName",
				UserAgent: "tstUserAgent",
				FileType:  "testFileType",
			},
			wantErr: false,
		},
		{
			name: "Invalid file name",
			args: args{
				fileName: "invalid file",
				db:       db,
			},
			want:    model.UploadFileLog{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := selectFileLog(tt.args.fileName, tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("selectFileLog() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("selectFileLog() got = %v, want %v", got, tt.want)
			}
		})
	}
}
