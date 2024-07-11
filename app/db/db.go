package db

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

var _db *gorm.DB

var (
	dbInterfaces = []interface{}{
		User{},
		Session{},
		Investment{},
		InvestmentPrice{},
		InvestmentTransaction{},
	}
)

func IsRecordNotFoundError(e error) bool {
	return jerr.HasError(e, "record not found")
}

func getDb() (*gorm.DB, error) {
	if _db == nil {
		var err error
		_db, err = gorm.Open("sqlite3", "money.db")
		if err != nil {
			return _db, jerr.Get("Failed to connect to database", err)
		}
		_db.LogMode(false)
		for _, iface := range dbInterfaces {
			result := _db.AutoMigrate(iface)
			if result.Error != nil {
				return result, result.Error
			}
		}
	}
	return _db, nil
}

func create(value interface{}) *gorm.DB {
	db, _ := getDb()
	result := db.Create(value)
	return result
}

func find(out interface{}, where ...interface{}) *gorm.DB {
	db, err := getDb()
	if err != nil {
		log.Fatalf("error getting db; %v", err)
	}
	result := db.Find(out, where...)
	return result
}

func save(value interface{}) *gorm.DB {
	db, _ := getDb()
	if db.Error != nil {
		return db
	}
	result := db.Save(value)
	return result
}

func remove(value interface{}) *gorm.DB {
	db, _ := getDb()
	if db.Error != nil {
		return db
	}
	result := db.Delete(value)
	return result
}
