package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// func CreateDSN(username, password, host, port, db string) string {
// 	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True", username, password, host, port, db)
// }

func NewDabataseConn(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

// "course:secret@tcp(db:3306)/course?charset=utf8mb4&parseTime=True&loc=Local"

// func NewDabataseConn() *gorm.DB {
// 	dsn := "root:root@tcp(127.0.0.1:60082)/course?charset=utf8mb4&parseTime=True&loc=Local"
// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		panic(err)
// 	}
// 	return db
// }
