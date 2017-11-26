package main

import (
	"fmt"
	"log"
	"time"
)

const (
	NUMCLIENTS = 12
	MAXCLIENTS = 5
	NUMBARBERS = 2
	CUT        = 0
	OUT        = 1
	WAIT       = 2
	GO         = 3
)

type Client struct {
	ch    chan int
	state int
}

type Recep struct {
	cut     chan Client
	clients chan Client
	exit    chan bool
}

type Barber struct {
	ch chan Client
}

var recept = Recep{make(chan Client), make(chan Client), make(chan bool)}
var barbers = Barber{make(chan Client)}
var done = make(chan string)

func printState(state int) {
	switch state {
	case CUT:
		fmt.Println(" CUT")
	case OUT:
		fmt.Println(" OUT")
	case WAIT:
		fmt.Println(" WAIT")
	}
}

func recepcionist() {
	wRoom := 0
	ext := false

	for !ext {
		select {
		case c := <-recept.clients:
			if wRoom == MAXCLIENTS {
				c.ch <- OUT
			} else {
				wRoom++
				c.ch <- WAIT
			}
			// Test limits...
			if wRoom > MAXCLIENTS {
				log.Fatal("Error max")
			}
		case <-recept.cut:
			wRoom--
			// Test limits...
			if wRoom < 0 {
				log.Fatal("Error min")
			}
		case ext = <-recept.exit:
			close(barbers.ch)
		}
	}
	done <- ""
}

func client(id int) {
	c := Client{make(chan int), GO}

	fmt.Println(fmt.Sprintf("Cliente : %d voy a cortarme el pelo", id))
	recept.clients <- c // Go to barbery...
	state := <-c.ch     // Wait response...
	if state == WAIT {
		c.state = CUT
		barbers.ch <- c
		state = <-c.ch
		fmt.Println(fmt.Sprintf("Cliente : %d me he cortado el pelo", id))
		if state != CUT {
			log.Fatal("Erro bad state of client.")
		}
	} else if state == OUT {
		fmt.Printf("Cliente : %d me voy de la barberia, esta llena\n", id)
	} else {
		log.Fatal("Error bad state of client...")
	}
	done <- ""
}

func barber(id int) {
	fmt.Printf("Barbero %d : me duermo esperando clientes\n", id)
	for b := range barbers.ch {
		fmt.Printf("Barbero %d : empiezo a cortar el pelo\n", id)
		time.Sleep(1000 * time.Millisecond)
		recept.cut <- b
		b.ch <- CUT
		fmt.Printf("Barbero %d : me duermo esperando clientes\n", id)
	}
}

func main() {
	//  Launch barbery...
	go recepcionist()
	for i := 0; i < NUMBARBERS; i++ {
		go barber(i)
	}
	// Launch Client...
	for i := 0; i < NUMCLIENTS; i++ {
		go client(i)
	}
	// Wait cut clients...
	for i := 0; i < NUMCLIENTS; i++ {
		<-done
	}
	recept.exit <- true
	<-done
}
