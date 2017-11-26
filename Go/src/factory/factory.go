package main

import (
	"fmt"
	"sem"
	"sync"
	"time"
)

const (
	MAXELEMENTS = 10
	MAXROBOTS   = 3
	NUMMOBILES  = 2
	PRODWIRES   = 30
	PRODCASES   = 6
	PRODSCREENS = 6
	PRODBOARDS  = 6
)

const (
	MOBWIRES   = 5
	MOBBOARDS  = 1
	MOBCASES   = 1
	MOBSCREENS = 1
)

type Mobile struct {
	wires  []int
	case_  int
	screen int
	board  int
}

type prodCons struct {
	id   int
	buf  []int
	prod *sem.Sem
	cons *sem.Sem
	mut  *sem.Sem
	j    int
	i    int
}

var mutRobot = sem.NewSem(1)

var wg = sync.WaitGroup{}

func initbuf(buf []int) {
	for i := range buf {
		buf[i] = -1
	}
}

func NewProdCons() *prodCons {
	ps := &prodCons{}
	ps.buf = make([]int, MAXELEMENTS)
	initbuf(ps.buf)
	ps.prod = sem.NewSem(MAXELEMENTS)
	ps.cons = sem.NewSem(0)
	ps.mut = sem.NewSem(1)
	ps.j = 0
	return ps
}

func (ps *prodCons) put(id int) {
	ps.prod.Down() // One element added to chain.
	ps.mut.Down()

	ps.buf[ps.i] = id
	ps.i = (ps.i + 1) % MAXELEMENTS

	ps.mut.Up()
	ps.cons.Up() // You can get one new element int chain.
}

func (ps *prodCons) get() int {
	ps.cons.Down() // One element less in the chain.
	ps.mut.Down()  // Protect data.

	e := ps.buf[ps.j]
	ps.buf[ps.j] = -1
	ps.j = (ps.j + 1) % MAXELEMENTS

	ps.mut.Up()  // Get out data protected.
	ps.prod.Up() // Producer can put a new element.

	return e
}

func printWires(wires []int) {
	fmt.Print("[")
	for _, w := range wires {
		fmt.Print(" ")
		fmt.Print(w)
		fmt.Print(" ")
	}
	fmt.Print("], ")
}

func prodBoards(boards *prodCons) {
	defer wg.Done()

	idBoards := 0
	for i := 0; i < PRODBOARDS; i++ {
		boards.put(idBoards)
		idBoards++
	}
}

func prodScreens(screens *prodCons) {
	defer wg.Done()

	idScreens := 0
	for i := 0; i < PRODSCREENS; i++ {
		screens.put(idScreens)
		idScreens++
	}
}

func prodCases(cases *prodCons) {
	defer wg.Done()

	idCases := 0
	for i := 0; i < PRODCASES; i++ {
		cases.put(idCases)
		idCases++
	}
}

func prodWires(wires *prodCons) {
	defer wg.Done()

	idWires := 0
	for i := 0; i < PRODWIRES; i++ {
		wires.put(idWires)
		idWires++
	}
}

func getMobile(boards, screens, cases, wires *prodCons) *Mobile {
	mobile := &Mobile{}

	mobile.board = boards.get()
	mobile.screen = screens.get()
	mobile.case_ = cases.get()

	mobile.wires = make([]int, MOBWIRES)
	for i := 0; i < MOBWIRES; i++ {
		mobile.wires[i] = wires.get()
	}
	return mobile
}

func printMobile(mobile *Mobile, state string) {
	fmt.Print(" cables : ")
	printWires(mobile.wires)
	fmt.Print(" pantalla : ", int(mobile.screen))
	fmt.Print(" carcasa : ", mobile.case_)
	fmt.Print(" placa : ", mobile.board)
	fmt.Println(state)
}

func printRobot(id int, mobile *Mobile, state string) {
	fmt.Print("robot : ", id)
	printMobile(mobile, state)
}

func robot(id int, boards, screens, cases, wires *prodCons) {
	defer wg.Done()

	for i := 0; i < NUMMOBILES; i++ {
		mobile := getMobile(boards, screens, cases, wires)

		mutRobot.Down()
		printRobot(id, mobile, " Iniciando...")
		mutRobot.Up()

		time.Sleep(200 * time.Millisecond)

		mutRobot.Down()
		printRobot(id, mobile, " Finalizando...")
		mutRobot.Up()
	}
}

func main() {
	boards := NewProdCons()
	screens := NewProdCons()
	cases := NewProdCons()
	wires := NewProdCons()

	wg.Add(4)
	go prodBoards(boards)
	go prodCases(cases)
	go prodWires(wires)
	go prodScreens(screens)

	for i := 0; i < MAXROBOTS; i++ {
		wg.Add(1)
		go robot(i, boards, screens, cases, wires)
	}
	wg.Wait()
}
