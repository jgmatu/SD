package main

/*
URL: https://github.com/mccoyst/myip/blob/master/myip.go
URL: http://changsijay.com/2013/07/28/golang-get-ip-address/
*/

import (
	"bufio"
	"log"
	"logiclog"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	PORTS    = 1024
	IPv4     = 255
	PORTIP   = 2
	FIELDSV4 = 4
	ARGPROG  = 1
	ARGPORT  = ARGPROG + 1
)

var mut = sync.Mutex{}

func getIp() string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func isInvalidArgs() bool {
	return len(os.Args) != ARGPORT
}

func isPort(port string) bool {
	i, err := strconv.Atoi(port)

	return i >= PORTS && err == nil
}

func isBadPort() bool {
	return isInvalidArgs() || !isPort(os.Args[1])
}

func isBadSocket(ip string) bool {
	return ip == "" || isBadPort()
}

func getId() string {
	ip := getIp()
	if isBadSocket(ip) {
		os.Stderr.WriteString("Bad arguments the port is invalid\n")
		os.Exit(1)
	}
	port := os.Args[1]
	return ip + ":" + port
}

func handleConn(c net.Conn, logf *logiclog.Log) {
	defer c.Close()

	mut.Lock()
	logf.ReceiveConn(c)
	mut.Unlock()
}

func connSock(id string, logf *logiclog.Log) {
	bind, err := net.Listen("tcp", id)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := bind.Accept()
		if err != nil {
			log.Print(err)
			continue // e.g., connection aborted.
		}
		go handleConn(conn, logf)
	}
}

func getFile() string {
	filename, err := os.Hostname()

	if err != nil {
		log.Fatal(err)
	}
	return filename + ".txt"
}

func isinvalidIp_v4(field string) bool {
	i, err := strconv.Atoi(field)

	return i > IPv4 || err != nil
}

func isIpV4(ip string) bool {
	fields := strings.Split(ip, ".")

	if len(fields) != FIELDSV4 {
		return false
	}
	for _, field := range fields {
		if isinvalidIp_v4(field) {
			return false
		}
	}
	return true
}

func hasId(id string) bool {
	fields := strings.Split(id, ":")

	if len(fields) != PORTIP {
		return false
	}
	ip := fields[0]
	port := fields[1]
	return isIpV4(ip) && isPort(port)
}

func getLog(logf *logiclog.Log, event string) {
	scanner := bufio.NewScanner(strings.NewReader(event))
	text, id := "", ""
	pos := 0

	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		if pos == 0 {
			id = scanner.Text()
		} else {
			text += " " + scanner.Text()
		}
		pos++
	}

	if hasId(id) {
		mut.Lock()
		logf.SendConn(id, text)
		mut.Unlock()
	} else {
		text = id + text
		mut.Lock()
		logf.Event(text)
		mut.Unlock()
	}
}

func readInput(logf *logiclog.Log) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		getLog(logf, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	id := getId()
	filename := getFile()
	logf := logiclog.NewLog(id, filename)

	go connSock(id, logf)
	readInput(logf)
}
