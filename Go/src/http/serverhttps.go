package main

import (
	"crypto/tls"
	"log"
	"net"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
)

func main() {
	log.SetFlags(log.Lshortfile)

	cer, err := tls.LoadX509KeyPair("cert.pem", "privkey.pem")
	if err != nil {
		log.Println(err)
		return
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	ln, err := tls.Listen("tcp", ":443", config)
	if err != nil {
		log.Println(err)
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

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

func getDBInfo(conn net.Conn) string {
	db := connect()
	result := query(conn, db)
	db.Close()
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	n, err := conn.Write([]byte("Hello secure world!!!\n"))
	if err != nil {
		log.Println(n, err)
		return
	}
	conn.Write([]byte("List of users:\n"))
	info := getDBInfo(conn)
	conn.Write([]byte(info))
}
