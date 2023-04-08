package respository

import (
	log2 "github.com/gallifreyCar/gin-file-server/log"
	"github.com/gallifreyCar/gin-file-server/model"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

func getDataBase() (db *gorm.DB, err error) {

	//set dataBase.log
	file := log2.InitLogFile("dataBase.log", "[DataBase]")
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

			log.Fatal(err)
		}
	}(file)

	// Capture connection properties.
	cfg := mysqlDriver.Config{
		User:   os.Getenv("DBUser"),
		Passwd: os.Getenv("DBPassword"),
		DBName: os.Getenv("DBName"),
		Addr:   "localhost:3306",
		Net:    "tcp",
		Params: map[string]string{
			"loc":       "Local",
			"parseTime": "True",
		},
	}
	//log cfg
	logCfg := mysqlDriver.Config{
		User:   "DBUser",
		Passwd: "DBPassword",
		DBName: os.Getenv("DBName"),
		Addr:   "localhost:3306",
		Net:    "tcp",
		Params: map[string]string{
			"loc":       "Local",
			"parseTime": "True",
		},
	}

	db, err = gorm.Open(mysql.Open(cfg.FormatDSN()), &gorm.Config{})
	if err != nil {
		log.Println(err)
	} else {
		log.Println(logCfg.FormatDSN())
	}
	return db, err
}
func insertFileLog(savePath, fileName, userAgent, fileType string, db *gorm.DB) (ID uint, RowsAffected int64, err error) {

	//set dataBase.log
	file := log2.InitLogFile("dataBase.log", "[InsertFile]")
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

			log.Fatal(err)
		}
	}(file)

	//create record insert into database
	fileLog := model.UploadFileLog{
		FileName:  fileName,
		UserAgent: userAgent,
		FileType:  fileType,
		SavePath:  savePath,
	}
	result := db.Create(&fileLog)
	if result.Error == nil {
		log.Printf("ID:%v,RowsAffected:%v", fileLog.ID, result.RowsAffected)
	}

	return fileLog.ID, result.RowsAffected, result.Error

}
func selectFileLog(fileName string, db *gorm.DB) (model.UploadFileLog, error) {

	//set dataBase.log
	file := log2.InitLogFile("dataBase.log", "[selectFile]")
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

			log.Fatal(err)
		}
	}(file)

	var fileLog model.UploadFileLog
	result := db.Where(&model.UploadFileLog{FileName: fileName}).Last(&fileLog)

	log.Printf("ID: %v,RowsAffected: %v", fileLog.ID, result.RowsAffected)

	return fileLog, result.Error
}
