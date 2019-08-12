package main

import (
	"fmt"
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
	// clientStore.Create()
	tokenStore, _ := mysql.NewTokenStore(db)

	manager.MapClientStorage(clientStore)
	manager.MapTokenStorage(tokenStore)

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

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			id := r.FormValue("id")
			secrete := r.FormValue("secrete")
			fmt.Println(id, secrete)
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		http.Error(w, "Only support post", http.StatusBadRequest)

	})

	http.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		err := srv.HandleAuthorizeRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		srv.HandleTokenRequest(w, r)
	})

	log.Fatal(http.ListenAndServe(":9096", nil))

}
