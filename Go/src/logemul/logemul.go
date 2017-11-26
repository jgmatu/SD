package main

import (
      "logiclog"
      "sync"
)

type Channels struct {
      AB chan logiclog.Msg
      BC chan logiclog.Msg
      CD chan logiclog.Msg
      DA chan logiclog.Msg
}

var wg = sync.WaitGroup{}

func NodeA(chans Channels , id , file string) {
      defer wg.Done()
      logf := logiclog.NewLog(id, file)

      logf.Event("Hey!")
      logf.Send("Dinner!", chans.AB)
      logf.Event("Cool!")
      logf.Event("Arg!")
      logf.Event("kk")
      logf.Receive(chans.DA)
      logf.Event("Exit")
      logf.Close()
}

func NodeB(chans Channels, id , file string) {
      defer wg.Done()
      logf := logiclog.NewLog(id, file)

      logf.Event("Ho!")
      logf.Event("Bored")
      logf.Receive(chans.AB)
      logf.Event("Good")
      logf.Send("Din!", chans.BC)
      logf.Event("Done!")
      logf.Close()
}

func NodeC(chans Channels, id , file string) {
      defer wg.Done()
      logf := logiclog.NewLog(id, file)

      logf.Event("Lets!")
      logf.Event("Bored!")
      logf.Event("Very Bored!")
      logf.Receive(chans.BC)
      logf.Send("Go to Dinner!", chans.CD)
      logf.Event("Wake!")
      logf.Event("Exit!")
      logf.Close()
}
func NodeD(chans Channels, id , file string) {
      defer wg.Done()
      logf := logiclog.NewLog(id, file)

      logf.Event("Go!")
      logf.Event("Wait")
      logf.Event("WaitDone!")
      logf.Receive(chans.CD)
      logf.Send("I Love You!" , chans.DA)
      logf.Event("Dress")
      logf.Event("Exit")
      logf.Close()
}

func NewChans() Channels {
      chans := Channels{}

      chans.AB = make(chan logiclog.Msg)
      chans.BC = make(chan logiclog.Msg)
      chans.CD = make(chan logiclog.Msg)
      chans.DA = make(chan logiclog.Msg)
      return chans
}

func main() {
      chans := NewChans()

      wg.Add(4)
      go NodeA(chans, "a" , "A.txt")
      go NodeB(chans, "b" , "B.txt")
      go NodeC(chans, "c" , "C.txt")
      go NodeD(chans, "d" , "D.txt")
      wg.Wait()
}
