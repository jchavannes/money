package db

import (
	"github.com/jchavannes/gorm"
	_ "github.com/jchavannes/gorm/dialects/mysql"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/money/data/config"
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

func isRecordNotFoundError(e error) bool {
	return e.Error() == "record not found"
}

func getDb() (*gorm.DB, error) {
	if _db == nil {
		conf := config.GetMysqlConfig()
		var err error
		connectionString := conf.Username + ":" + conf.Password + "@tcp(" + conf.Host + ")/" + conf.Database + "?parseTime=true"
		_db, err = gorm.Open("mysql", connectionString)
		_db.LogMode(false)
		if err != nil {
			return _db, jerr.Get("Failed to connect to database", err)
		}
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
	db, _ := getDb()
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
