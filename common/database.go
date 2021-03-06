package common

import (
	"fmt"
	"go-backend/entity"
	"net/url"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	driverName := viper.GetString("datasource.driverName")
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	loc := viper.GetString("datasource.loc")
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
		username,
		password,
		host,
		port,
		database,
		charset,
		url.QueryEscape(loc))
	// fmt.Printf(args)
	db, err := gorm.Open(driverName, args)
	// db 使用 golang 原生 sql.Open 自带连接池
	if err != nil {
		panic("fail to connect database, err: " + err.Error())
	}
	db.AutoMigrate(&entity.User{})
	db.AutoMigrate(&entity.Company{})
	db.AutoMigrate(&entity.CompanyUser{})
	db.AutoMigrate(&entity.Biology{})
	db.AutoMigrate(&entity.BiologyType{})
	db.AutoMigrate(&entity.FixedDevice{})
	db.AutoMigrate(&entity.PortableDevice{})
	db.AutoMigrate(&entity.FixedDeviceType{})
	db.AutoMigrate(&entity.PortableDeviceType{})
	db.AutoMigrate(&entity.FenceRecord{})
	db.AutoMigrate(&entity.EpidemicPrevention{})
	db.AutoMigrate(&entity.MedicalHistory{})
	db.AutoMigrate(&entity.OperationHistory{})
	
	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}
