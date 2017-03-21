package main

import (
	"fmt"
	"log"
)

const (
	NUMCLIENTS = 12
	MAXCLIENTS = 5
	NUMBARBERS = 2
	CUT        = 1
	OUT        = 2
)

type Client struct {
	ch chan int
}

type Recep struct {
	cut     chan Client
	clients chan Client
	exit    chan bool
}

type Barber struct {
	ch chan Client
}

type Barbery struct {
	recept  *Recep
	barbers *Barber
	done    chan string
	enable  bool
}

type PubBarber interface {
	NewBarbery() *Barbery
	RunBarbery()
	Barber(id int)
	Client(id int)
	WaitDone()
	Exit()
}

func (bb *Barbery) recepcionist() {
	wRoom := 0
	ext := false
	recept := bb.recept
	barbers := bb.barbers

	for !ext {
		select {
		case c := <-recept.clients:

			if wRoom == MAXCLIENTS || NUMBARBERS <= 0 {
				c.ch <- OUT
			} else {
				wRoom++
				c.ch <- CUT
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
			for i := 0; i < NUMBARBERS; i++ {
				bb.WaitDone()
			}
		}
	}
	bb.done <- ""
}

func (bb *Barbery) client(id int) {
	c := Client{make(chan int)}
	recept := bb.recept
	barbers := bb.barbers

	recept.clients <- c // Go to barbery...
	state := <-c.ch     // Wait response...

	if state == OUT {

		fmt.Println(fmt.Sprintf("Cliente %d : me voy de la barberia, esta llena", id))

	} else if state == CUT {

		fmt.Printf("Cliente %d : me siento en la sala de espera\n", id)
		barbers.ch <- c
		fmt.Printf("Cliente %d : me corto el pelo\n", id)
		<-c.ch
		fmt.Printf("Cliente %d : termino de cortarme el pelo\n", id)

	} else {

		log.Fatal("Error state of client")

	}
	bb.done <- ""
}

func (bb *Barbery) barber(id int) {
	fmt.Printf("Barbero %d : me duermo esperando clientes\n", id)
	for b := range bb.barbers.ch {
		bb.recept.cut <- b
		fmt.Printf("Barbero %d : empiezo a cortar el pelo\n", id)
		b.ch <- CUT
		fmt.Printf("Barbero %d : termino de cortar el pelo\n", id)
		fmt.Printf("Barbero %d : me duermo esperando clientes\n", id)
	}
	bb.done <- ""
}

func NewBarbery() *Barbery {
	bb := &Barbery{}

	bb.recept = &Recep{}
	bb.recept.cut = make(chan Client)
	bb.recept.clients = make(chan Client)
	bb.recept.exit = make(chan bool)

	bb.barbers = &Barber{}
	bb.barbers.ch = make(chan Client)

	bb.done = make(chan string)

	bb.enable = false

	return bb
}

func (bb *Barbery) RunBarbery() {
	if bb.enable {
		log.Fatal("Barbery run...")
		return
	}
	go bb.recepcionist()
	bb.enable = true
}

func (bb *Barbery) Barber(id int) {
	go bb.barber(id)
}

func (bb *Barbery) Client(id int) {
	go bb.client(id)
}

func (bb *Barbery) WaitDone() {
	<-bb.done
}

func (bb *Barbery) Exit() {
	bb.recept.exit <- true
}

func main() {
	bb := NewBarbery()

	bb.RunBarbery()
	for i := 0; i < NUMBARBERS; i++ {
		bb.Barber(i)
	}
	// Launch Client...
	for i := 0; i < NUMCLIENTS; i++ {
		bb.Client(i)
	}
	// Wait cut clients...
	for i := 0; i < NUMCLIENTS; i++ {
		bb.WaitDone()
	}
	bb.Exit()
	bb.WaitDone()
}
