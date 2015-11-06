package main

import (
  "os"
)

func main() {
  hostName, _ := os.Hostname()
  println("Hello from", hostName)
}
