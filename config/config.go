package config

import (
	"alta/be4/mvc/models"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	//konfigurasi koneksi DB
	// config := map[string]string{
	// 	"DB_Username": "root",
	// 	"DB_Password": "qwerty123",
	// 	"DB_Port":     "3306",
	// 	"DB_Host":     "127.0.0.1",
	// 	"DB_Name":     "be4_mvc",
	// }
	// connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
	// 	config["DB_Username"],
	// 	config["DB_Password"],
	// 	config["DB_Host"],
	// 	config["DB_Port"],
	// 	config["DB_Name"])
	// connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
	// 	os.Getenv("DB_USERNAME"),
	// 	os.Getenv("DB_PASSWORD"),
	// 	config["DB_Host"],
	// 	config["DB_Port"],
	// 	config["DB_Name"])

	connectionString := os.Getenv("CONNECTION_STRING")
	//connect ke DB
	var err error
	DB, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	//menjalankan Auto Migrate
	InitiateMigrate()
}

func InitiateMigrate() {
	DB.AutoMigrate(&models.User{})
	// DB.AutoMigrate(&models.Book{})
	// DB.AutoMigrate(&models.ProductType{})
	// DB.AutoMigrate(&models.Operator{})
	// DB.AutoMigrate(&models.Product{})
	// DB.AutoMigrate(&models.ProductDescription{})

}

func InitDBTest() {
	connectionStringTest := "root:qwerty123@tcp(127.0.0.1:3306)/be4_mvc_test?charset=utf8&parseTime=True&loc=Local"
	//connect ke DB
	var err error
	DB, err = gorm.Open(mysql.Open(connectionStringTest), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	//menjalankan Auto Migrate
	InitiateMigrateTest()
}

func InitiateMigrateTest() {
	DB.Migrator().DropTable(&models.User{})
	DB.AutoMigrate(&models.User{})
}
