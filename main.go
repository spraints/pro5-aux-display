package main

import (
  "flag"
  "fmt"
)

func main() {
  var listenPort = flag.Int("port", 3000, "The port that this server listens on.")
  var remoteHost = flag.String("pro5-host", "localhost", "The name of the pro5 server.")
  var remotePort = flag.Int("pro5-port", 9002, "The port that pro5's remote stage display is listening on.")
  var remotePassword = flag.String("pro5-password", "stage", "The password to the pro5 remote stage display")

  flag.Parse()

  fmt.Printf("listen on %d\n", *listenPort)
  fmt.Printf("connect to %s:%d\n", *remoteHost, *remotePort)
  fmt.Printf("password = %s\n", *remotePassword)
}
