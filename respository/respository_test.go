package respository

import "testing"

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

	id, rowsAffected, err := insertFileLog("testSavePath", "testFileName", "tstUserAgent", "testFileType")
	t.Logf("ID: %v,RowsAffected: %v", id, rowsAffected)
	if err != nil {
		t.Error(err)
	}

}

func TestSelectFileLog(t *testing.T) {

	fileLog, err := selectFileLog("testFileName")
	t.Logf("fileLog: %v", fileLog)
	if err != nil {
		t.Error(err)
	}

}
