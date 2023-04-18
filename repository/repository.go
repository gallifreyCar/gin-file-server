package repository

import (
	"github.com/gallifreyCar/gin-file-server/m-logger"
	"github.com/gallifreyCar/gin-file-server/model"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

func GetDataBase() (db *gorm.DB, err error) {

	//set database zap logger
	logger, err, closeFunc := m_logger.InitZapLogger("gin-file-server.log", "[GetDataBase]")
	logger.Error("Fail to init zap logger", zap.Error(err))
	defer closeFunc()

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

	//set a zap logger
	logger, err, closeFunc := m_logger.InitZapLogger("gin-file-server.log", "[InsertFileLog]")
	logger.Error("Fail to init zap logger", zap.Error(err))
	defer closeFunc()

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
		logger.Error("Fail to create file log in repository", zap.Error(result.Error))
	}

	logger.Info("Insert file log success!", zap.Uint("fileLog.ID", fileLog.ID), zap.Int64("row affected", result.RowsAffected))
	return fileLog.ID, result.RowsAffected, result.Error

}
func SelectFileLog(fileName string, db *gorm.DB) (model.UploadFileLog, error) {

	//set a zap logger
	logger, err, closeFunc := m_logger.InitZapLogger("gin-file-server.log", "[SelectFileLog]")
	logger.Error("Fail to init zap logger", zap.Error(err))
	defer closeFunc()

	var fileLog model.UploadFileLog
	result := db.Where(&model.UploadFileLog{FileName: fileName}).Last(&fileLog)

	if result.Error != nil {
		logger.Error("File to use gorm select file log by file name", zap.Error(result.Error))
		return fileLog, result.Error
	}

	logger.Info("Select file by file name success!", zap.Uint("fileLog.ID", fileLog.ID))
	return fileLog, result.Error
}
