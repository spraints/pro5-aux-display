package main

import (
  "flag"

  "pro5state"
  "pro5stage"
  "pro5web"
)

func main() {
  connectInfo := pro5stage.ConnectInfo{}
  serverInfo := pro5web.WebServerInfo{}

  flag.IntVar(&serverInfo.Port, "port", 3000, "The port that this server listens on.")
  flag.StringVar(&connectInfo.Host, "pro5-host", "localhost", "The name of the pro5 server.")
  flag.IntVar(&connectInfo.Port, "pro5-port", 9002, "The port that pro5's remote stage display is listening on.")
  flag.StringVar(&connectInfo.Password, "pro5-password", "stage", "The password to the pro5 remote stage display")

  flag.Parse()

  state := pro5state.New()
  go pro5stage.Run(connectInfo, state)
  pro5web.StartServer(serverInfo, state)
}
