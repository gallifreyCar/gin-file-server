package model

import (
	"gorm.io/gorm"
)

/*
type uploadFileLog struct {
	ID        uint           `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	fileName    string `gorm:"->;<-:create"` //allow read and write(create)
	userAgent   string `gorm:"->;<-:create"` //allow read and write(create)
	fileType string `gorm:"->;<-:create"` //allow read and write(create)
	SavePath  string `gorm:"->;<-:create"` //allow read and write(create)
}
*/

// UploadFileLog equals upper uploadFileLog struct
type UploadFileLog struct {
	gorm.Model
	FileName  string `gorm:"->;<-:create"` //allow read and write(create)
	UserAgent string `gorm:"->;<-:create"` //allow read and write(create)
	FileType  string `gorm:"->;<-:create"` //allow read and write(create)
	FileSize  string `gorm:"->;<-:create"` //allow read and write(create)
	SavePath  string `gorm:"->;<-:create"` //allow read and write(create)
}
