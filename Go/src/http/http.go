package main

import (
	"fmt"
	"log"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
)

func connect() *sql.DB {
	db, err := sql.Open("mysql", "javi:c1br4n13@tcp(mydb0.clbuzztihhqs.eu-central-1.rds.amazonaws.com:3306)/data0")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func query(db *sql.DB) string {
        id := 0
        name := ""
        apellido := ""
        result := ""

        rows, err := db.Query("Select * from users")
        if err != nil {
                log.Fatal(err)
        }
        for rows.Next() {
                err := rows.Scan(&id, &name, &apellido)
                if err != nil {
                        log.Fatal(err)
                }
                result += fmt.Sprintf("%d : %s\t%s\n", id , name , apellido)
        }
        err = rows.Err()
        if err != nil {
                log.Fatal(err)
        }
        return result
}

func getDBInfo() string {
	db := connect()
	result := query(db)
	db.Close()
	return result
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love this site!\n")
	info := getDBInfo()
	fmt.Fprintf(w , "\nTable of users : \n%s\n", info)
}

func main() {
    http.HandleFunc("/", handler)
    err := http.ListenAndServe(":80", nil)
    if err != nil {
            log.Fatal(err)
    }
}
