package main

import (
  "flag"
  "log"

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

  var pro5, err = pro5stage.ConnectToPro5(connectInfo)
  if err != nil {
    log.Fatal("ConnectToPro5: ", err)
  }
  state := pro5state.New()
  go pro5.ReadEverything(state)
  pro5web.StartServer(serverInfo, state)
}
