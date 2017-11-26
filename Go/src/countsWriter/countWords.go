package main

import (
      "fmt"
      "bufio"
      "bytes"
)

type CountWords int

func (c *CountWords) Write(p []byte) (int ,error) {
      scanner := bufio.NewScanner(bytes.NewReader(p))
      scanner.Split(bufio.ScanWords)
      nWords := 0
      for scanner.Scan() {
            nWords++
      }
      *c += CountWords(nWords)
      err := scanner.Err()
      return nWords, err
}


func main() {
      var c CountWords
      fmt.Fprintln(&c , "Hello world!!")
      fmt.Println(c)
}
