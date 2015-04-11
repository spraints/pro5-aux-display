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

type MessageStream interface {
  Tap(messages chan<- string)
}

func StartServer(info WebServerInfo, stream MessageStream) {
  var mux = http.NewServeMux()
  mux.Handle("/connect", serveWebSockets(stream))
  mux.Handle("/", http.FileServer(http.Dir("public/")))
  http.ListenAndServe(fmt.Sprintf(":%d", info.Port), mux)
  // This never returns.
}

func serveWebSockets(stream MessageStream) websocket.Handler {
  return func(ws *websocket.Conn) {
    messages := make(chan string, 5)
    stream.Tap(messages)
    for message := range messages {
      fmt.Fprintf(ws, "%s", message)
    }
  }
}
