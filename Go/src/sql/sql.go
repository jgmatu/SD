package main

import (
        "database/sql"
        _ "github.com/go-sql-driver/mysql"
        "log"
        "fmt"
)

func connect() *sql.DB {
      db, err := sql.Open("mysql", "javi:c1br4n13@tcp(mydb.clbuzztihhqs.eu-central-1.rds.amazonaws.com:3306)/data0")
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

func main() {
      db := connect()
      info := query(db)
      fmt.Println(info)
      db.Close()
}
