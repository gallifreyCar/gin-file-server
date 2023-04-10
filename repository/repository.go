package repository

import (
	log2 "github.com/gallifreyCar/gin-file-server/log"
	"github.com/gallifreyCar/gin-file-server/model"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

func GetDataBase() (db *gorm.DB, err error) {

	//set dataBase.log
	file, log := log2.InitLogFile("gin-file-server.log", "[DataBase]")
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
func InsertFileLog(savePath, fileName, userAgent, fileType string, fileSize int64, db *gorm.DB) (ID uint, RowsAffected int64, err error) {

	//set dataBase.logger
	file, logger := log2.InitLogFile("gin-file-server.log", "[InsertFile]")
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

			logger.Fatal(err)
		}
	}(file)

	//create record insert into database
	fileLog := model.UploadFileLog{
		FileName:  fileName,
		UserAgent: userAgent,
		FileType:  fileType,
		SavePath:  savePath,
		FileSize:  fileSize,
	}
	result := db.Create(&fileLog)
	if result.Error == nil {
		logger.Printf("ID:%v,RowsAffected:%v", fileLog.ID, result.RowsAffected)
	}

	return fileLog.ID, result.RowsAffected, result.Error

}
func SelectFileLog(fileName string, db *gorm.DB) (model.UploadFileLog, error) {

	//set dataBase.logger
	file, logger := log2.InitLogFile("gin-file-server.log", "[selectFile]")
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

			logger.Fatal(err)
		}
	}(file)

	var fileLog model.UploadFileLog
	result := db.Where(&model.UploadFileLog{FileName: fileName}).Last(&fileLog)

	if result.Error != nil {
		logger.Println(result.Error)
		return fileLog, result.Error
	}

	logger.Printf("ID: %v,RowsAffected: %v", fileLog.ID, result.RowsAffected)
	return fileLog, result.Error
}
