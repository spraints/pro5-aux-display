package pro5web

import (
  "fmt"
  "net/http"

  "golang.org/x/net/websocket"
)

// Arguments for serving the web front end.
type WebServerInfo struct {
  Port int
}

type NewChannels chan<- (chan<- string)

func StartServer(info WebServerInfo, messageChannels NewChannels) {
  var mux = http.NewServeMux()
  mux.Handle("/connect", serveWebSockets(messageChannels))
  mux.Handle("/", http.FileServer(http.Dir("public/")))
  http.ListenAndServe(fmt.Sprintf(":%d", info.Port), mux)
  // This never returns.
}

func serveWebSockets(messageChannels NewChannels) websocket.Handler {
  return func(ws *websocket.Conn) {
    messages := make(chan string, 5)
    messageChannels <- messages
    for message := range messages {
      fmt.Fprintf(ws, "%s", message)
    }
  }
}
