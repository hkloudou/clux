package clux

import "gorm.io/gorm"

var db *gorm.DB

func GetDb() *gorm.DB {
	if _verbose {
		return db.Debug()
	}
	return db
}
