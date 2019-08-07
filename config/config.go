package config

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	mysql "github.com/imrenagi/go-oauth2-mysql"
	"github.com/jmoiron/sqlx"
	"gopkg.in/oauth2.v3/models"
)

// CreateUser ddd
func CreateUser() {
	db, err := sqlx.Connect("mysql", "root:abc123@tcp(127.0.0.1:3306)/oauth2db?parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}

	// manager := manage.NewDefaultManager()

	clientStore, _ := mysql.NewClientStore(db, mysql.WithClientStoreTableName("oauth2db"))

	clientStore.Create(&models.Client{
		ID:     "222222",
		Secret: "22222222",
		Domain: "http://localhost:9094"})
	// tokenStore, _ := mysql.NewTokenStore(db)

}
