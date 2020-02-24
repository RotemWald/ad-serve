package main

import "github.com/jinzhu/gorm"

func NewDatabaseConnection(dialect string, databaseName string) *gorm.DB {
	db, err := gorm.Open(dialect, databaseName)
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
