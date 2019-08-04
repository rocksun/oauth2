package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	mysql "github.com/imrenagi/go-oauth2-mysql"
	"github.com/jmoiron/sqlx"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
)

func main() {
	db, err := sqlx.Connect("mysql", "root:abc123@tcp(127.0.0.1:3306)/oauth2db?parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}

	manager := manage.NewDefaultManager()

	clientStore, _ := mysql.NewClientStore(db, mysql.WithClientStoreTableName("oauth2db"))
	tokenStore, _ := mysql.NewTokenStore(db)

	manager.MapTokenStorage(tokenStore)
	manager.MapClientStorage(clientStore)
	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	log.Fatal(http.ListenAndServe(":9096", nil))
}
