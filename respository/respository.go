package respository

import (
	"github.com/gallifreyCar/gin-file-server/model"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

func getDataBase() (db *gorm.DB, err error) {
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
	db, err = gorm.Open(mysql.Open(cfg.FormatDSN()), &gorm.Config{})

	return db, err
}
func insertFileLog(savePath, fileName, userAgent, fileType string, db *gorm.DB) (ID uint, RowsAffected int64, err error) {
	//set dataBase.log
	file, err := os.OpenFile("../log/dataBase.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)

	}

	log.SetPrefix("[InsertFile] ")
	log.SetOutput(file)
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
	file, err := os.OpenFile("../log/dataBase.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)

	}

	log.SetPrefix("[selectFile] ")
	log.SetOutput(file)
	defer func(file *os.File) {

		err := file.Close()
		if err != nil {

			log.Fatal(err)
		}
	}(file)

	var fileLog model.UploadFileLog
	//result := db.Where(&model.UploadFileLog{FileName: fileName}).First(fileLog)
	result := db.Where(&model.UploadFileLog{FileName: fileName}).First(&fileLog)

	log.Print(result.Error)
	log.Print(result.RowsAffected)

	return fileLog, result.Error
}
