package main

import (
      "fmt"
)

func main() {
      naturals := make(chan int)
      squares := make(chan int)

      // Counter
      go func() {
            for x := 0 ; x < 10 ; x++ {
                  naturals <- x
            }
            close(naturals)
      }()

      // Squarer
      go func() {
            for x := range naturals {
                  squares <- x * x
            }
            close(squares)
      }()

      // Printer (in main go routine)
      for y := range squares {
            fmt.Println(y)
      }
}
