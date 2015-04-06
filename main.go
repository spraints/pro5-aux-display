package main

import (
  "flag"
  "fmt"
)

type Pro5ConnectInfo struct {
  Host string
  Port int
  Password string
}

func main() {
  connectInfo := Pro5ConnectInfo{}

  var listenPort = flag.Int("port", 3000, "The port that this server listens on.")
  flag.StringVar(&connectInfo.Host, "pro5-host", "localhost", "The name of the pro5 server.")
  flag.IntVar(&connectInfo.Port, "pro5-port", 9002, "The port that pro5's remote stage display is listening on.")
  flag.StringVar(&connectInfo.Password, "pro5-password", "stage", "The password to the pro5 remote stage display")

  flag.Parse()

  fmt.Printf("listen on %d\n", *listenPort)

  ConnectToPro5(connectInfo)
}

type Pro5Connection struct {
}

func ConnectToPro5(info Pro5ConnectInfo) Pro5Connection {
  fmt.Printf("connect to %s:%d with pw=%s\n", info.Host, info.Port, info.Password)
  return Pro5Connection{}
}
