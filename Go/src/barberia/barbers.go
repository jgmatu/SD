package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

const (
	MAXCLIENTS = 5
	MAXBARBERS = 2
	NUMCLIENTS = 13
)

type Client struct {
	id int
	ch chan string
}

var RWMutex = sync.RWMutex{}

var reception = make(chan Client)
var barbers = make(chan Client)
var cut = make(chan Client)
var exit = make(chan bool)
var clients = make(chan string)
var done = make(chan string)

func (c Client) String() string {
	client := fmt.Sprintf("Client \n")
	client += fmt.Sprintf("---------- \n")
	client += fmt.Sprintf("Client : %d \n", c.id)
	client += fmt.Sprintf("Channel : %v \n", c.ch)
	client += fmt.Sprintf("---------- \n")
	return client
}

func printClients(clients []Client) {
	for _, c := range clients {
		fmt.Println(c)
	}
}

func initClients(max int) []Client {
	clients := make([]Client, max)

	for i := range clients {
		clients[i] = errClient()
	}
	return clients
}

func nClients(clients []Client) int {
	n := 0
	for _, c := range clients {
		if c.id >= 0 {
			n++
		}
	}
	return n
}

func isFull(clients []Client, max int) bool {
	return nClients(clients) >= max
}

func put(clients []Client, c Client, max int) bool {
	pos := 0
	add := false
	for pos < max && !add {
		if clients[pos].id < 0 {
			clients[pos] = c
			add = true
		} else {
			pos++
		}
	}
	return add
}

func errClient() Client {
	c := Client{}
	c.id = -1
	return c
}

func getC(client Client) Client {
	c := Client{}
	c.id = client.id
	c.ch = client.ch
	return c
}

func delete(clients []Client, c Client, max int) bool {
	pos := 0
	del := false
	for pos < max && !del {
		if clients[pos].id == c.id {
			clients[pos] = errClient()
			del = true
		} else {
			pos++
		}
	}
	return del
}

func get(clients []Client, max int) Client {
	c := errClient()
	pos := 0
	for pos < max && c.id < 0 {
		if clients[pos].id >= 0 {
			c = getC(clients[pos])
			clients[pos] = errClient()
		} else {
			pos++
		}
	}
	return c
}

func barber() {
	for c := range barbers {
		time.Sleep(200 * time.Millisecond)
		cut <- c
	}
}

func receptionist() {
	wRoom := initClients(MAXCLIENTS)
	armChBarb := initClients(MAXBARBERS)

	ext := false
	for !ext {
		select {
		case c := <-reception:
			if isFull(wRoom, MAXCLIENTS) {
				// Get out barbery...
				fmt.Println(fmt.Sprintf("Cliente : %d me voy de la barberia, esta llena", c.id))
				c.ch <- "Full"
			} else if isFull(armChBarb, MAXBARBERS) {
				// Go wating room...
				fmt.Println(fmt.Sprintf("Cliente : %d me siento en la sala de espera", c.id))
				if !put(wRoom, c, MAXCLIENTS) {
					log.Fatal("NOT PUT! IN WR")
				}
			} else {
				// Go to barber...
				if !put(armChBarb, c, MAXBARBERS) {
					log.Fatal("NOT PUT! DIR!!")
				}
				fmt.Println(fmt.Sprintf("Cliente : %d me corto el pelo", c.id))
				barbers <- c
			}
		case cc := <-cut:

			// Client cuted...
			fmt.Println(fmt.Sprintf("Cliente : %d Termino de cortarme el pelo", cc.id))
			cc.ch <- "Cut"

			// delete from armChBarb...
			if !delete(armChBarb, cc, MAXBARBERS) {
				log.Fatal("NOT DELETE!!! FROM ARM")
			}

			// Get the next from waitingRoom
			c := get(wRoom, MAXCLIENTS)
			if c.id >= 0 {
				if !put(armChBarb, c, MAXBARBERS) {
					log.Fatal("NOT PUT!! FROM WR")
				}
				fmt.Println(fmt.Sprintf("Cliente : %d me corto el pelo", c.id))
				barbers <- c
			}
		case ext = <-exit:
			fmt.Println("Closing barbery...")
			close(barbers)
		}
	}
	done <- "Closed"
}

func client(id int) {
	c := Client{id, make(chan string)}

	// I'm gonna cut my hair...
	reception <- c

	// Wait...
	<-c.ch

	clients <- "Finishing..."
}

func main() {
	RWMutex.RLock()
        RWMutex.RUnlock()
        RWMutex.Lock()
        RWMutex.Unlock()
	// Launch barbery...
	go receptionist()
	for i := 0; i < MAXBARBERS; i++ {
		go barber()
	}

	// Clients go to barbery...
	for i := 0; i < NUMCLIENTS; i++ {
		time.Sleep(100 * time.Millisecond)
		go client(i)
	}

	// Wait all clients cut...
	for i := 0; i < NUMCLIENTS; i++ {
		<-clients
	}
	exit <- true
	<-done
}
