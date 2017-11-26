package main

import (
      "fmt"
      "bufio"
      "bytes"
      "os"
)

type CountLines int

func (c *CountLines) Write(p []byte) (int, error) {
      scanner := bufio.NewScanner(bytes.NewReader(p))
      nlines := 0
      for scanner.Scan() {
            nlines++
      }
      *c += CountLines(nlines)
      err := scanner.Err()
      return nlines , err
}

func main() {
      var c CountLines

      fmt.Fprintf(os.Stdout , "Hello World!\n")
      fmt.Fprintf(&c , "Hello\n worlD\n\n\n")
      fmt.Println(c)
}
