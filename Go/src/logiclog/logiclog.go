package logiclog

import (
	"bufio"
	"fmt"
	"log"
	"logicclock"
	"net"
	"os"
	"strings"
)

type Log struct {
	id    string
	clock *logicclock.Clock
	file  *os.File
}

type Msg struct {
	id    string
	event string
	mark  string
}

type Line struct {
	mark *logicclock.Clock
	text string
}

// Creacion de un nuevo log.
func NewLog(id, filename string) *Log {
	logf := &Log{}
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)

	if err != nil {
		log.Fatalf("Error open file %s error : %v", filename, err)
	}
	logf.id = id
	logf.file = file
	logf.clock = logicclock.NewClock()
	return logf
}

// Escribir una linea de log en el fichero.
func (logf *Log) Event(event string) {
	logf.clock.Add(logf.id)
	s := fmt.Sprintf("%s ", logf.clock.Json()) + event + "\n"

	_, err := logf.file.WriteString(s)
	if err != nil {
		log.Fatalf("Error writing log : %s : error : %s", s, err)
	}
}

// Dep Message...
func (msg *Msg) String() string {
	s := fmt.Sprintf("--- Message --- \n")

	s += fmt.Sprintf("Id : %s\n", msg.id)
	s += fmt.Sprintf("Event : %s\n", msg.event)
	s += fmt.Sprintf("Mark : %s\n", msg.mark)
	s += fmt.Sprintf("------------------------\n")
	return s
}

func (logf *Log) putMsgRcv(mark, event string) {
	clock := logicclock.NewClock()

	clock.Data(mark)
	logf.clock.Max(clock) // Put our clock with the max value between two marks.
	logf.Event("Message Received : " + event)
}

// Recibo de mensajes con el reloj vectorial.
// Recibo el evento mas el.clock...
func (logf *Log) Receive(ch chan Msg) {
	msg := <-ch
	logf.putMsgRcv(msg.mark, msg.event)
}

// Envio de mensajes con el reloj vectorial.
// Envio el evento mas el.clock...
func (logf *Log) Send(event string, ch chan Msg) {
	logf.Event("Message Send : " + event)

	// Write in channel message to goroutine.
	msg := Msg{logf.id, event, logf.clock.Json()}
	ch <- msg
}

// Envio de mensajes por conexion de red.
func (logf *Log) SendConn(id, event string) {
	conn, err := net.Dial("tcp", id)
	if err != nil {
		log.Print(err)
		return
	}
	defer conn.Close()

	// Put message in log...
	logf.Event("Message Send : " + event)

	// Write in socket the message to peer.
	msg := string(logf.clock.Json()) + "\t" + event + "\n"
	fmt.Fprintf(conn, msg)
}

// Recibo de mensajes por conexion de red.
func (logf *Log) ReceiveConn(c net.Conn) {
	msg, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		log.Print(err) // e.g error receiving data from socket.
		return         // close connection with this client.
	}

	fields := strings.Split(msg, "\t")
	mark := fields[0]
	event := strings.Join(fields[1:], "\t")
	logf.putMsgRcv(mark, event[:len(event)-1]) // Drop a '\n' in event
}

// Close log.
func (logf *Log) Close() {
	logf.clock = nil
	logf.id = ""
	if err := logf.file.Close(); err != nil {
		log.Fatalf("Error closing file log : %v", err)
	}
}

/*
 ****************************************
 * 	Ordenacion de los logs		*
 ****************************************
 */

// Debug Line...
func (line *Line) String() string {
	var s string

	s += fmt.Sprintf("%s", line.mark.Json())
	s += fmt.Sprintf("%s\n", line.text)
	return s
}

// Obtain line following the format line in log.
func getLine(line string) Line {
	newLine := Line{}
	text := ""
	words := strings.Split(line, " ")

	newLine.mark = logicclock.NewClock()
	newLine.mark.Data(words[0])
	for i := 1; i < len(words); i++ {
		text += " " + words[i]
	}
	newLine.text = text
	return newLine
}

// Read log to compare casuals marks of lines all system logs...
func readlines(filename string) []Line {
	lines := make([]Line, 0)

	// Open file log...
	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatalf("Error open file %s error : %v\n", filename, err)
	}
	// Read File and get lines of log...
	scan := bufio.NewScanner(file)
	for scan.Scan() {
		lines = append(lines, getLine(scan.Text()))
	}
	if err := scan.Err(); err != nil {
		log.Fatalf("Error reading log : %v\n", err)
	}
	// Close file log...
	if err := file.Close(); err != nil {
		log.Fatalf("Error closing log : %v\n", err)
	}
	return lines
}

func getLogs(files []string) map[string][]Line {
	logs := make(map[string][]Line)

	for _, file := range files {
		logs[file] = readlines(file)
	}
	return logs
}

func getLines(log []Line) []Line {
	lines := make([]Line, 0)

	for _, line := range log {
		lines = append(lines, line)
	}
	return lines
}

// Obtain all the lines of all logs to order...
func getLinesLogs(logs map[string][]Line) []Line {
	lines := make([]Line, 0)

	for _, log := range logs {
		lines = append(lines, getLines(log)...)
	}
	return lines
}

func searchmin(lines []Line, from int) int {
	min := lines[from]
	posmin := from

	for i := from + 1; i < len(lines); i++ {
		less, err := lines[i].mark.IsLess(min.mark)
		if err == nil && less {
			min = lines[i]
			posmin = i
		}
	}
	return posmin
}

func sortByMarks(lines []Line) {
	for i := range lines {
		posmin := searchmin(lines, i)
		lines[i], lines[posmin] = lines[posmin], lines[i]
	}
}

func printLines(lines []Line) {
	for _, line := range lines {
		fmt.Printf("%s", line.String())
	}
}

func writeLines(filename string, lines []Line) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Error open file : %s err : %v", filename, err)
	}

	for _, line := range lines {
		file.WriteString(line.String())
	}
	file.Close()
}

func Order(filename string, files []string, output bool) {
	logs := getLogs(files)
	lines := getLinesLogs(logs)

	sortByMarks(lines)
	if output {
		printLines(lines)
	}
	writeLines(filename, lines)
}
