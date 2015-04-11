package pro5web

import (
  "fmt"
  "net/http"

  "golang.org/x/net/websocket"

  "pro5stage"
)

// Arguments for serving the web front end.
type WebServerInfo struct {
  Port int
}

func StartServer(info WebServerInfo, pro5 *pro5stage.Conn) {
  var mux = http.NewServeMux()
  mux.Handle("/connect", websocket.Handler(func(ws *websocket.Conn) { pro5.AddListener(ws) }))
  mux.Handle("/", http.FileServer(http.Dir("public/")))
  http.ListenAndServe(fmt.Sprintf(":%d", info.Port), mux)
}
