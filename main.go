package main

import (
  "container/list"
  "encoding/xml"
  "flag"
  "fmt"
  "io"
  "log"
  "net"
  "net/http"

  "golang.org/x/net/websocket"
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

  var pro5, err = ConnectToPro5(connectInfo)
  if err != nil {
    log.Fatal("ConnectToPro5: ", err)
  }
  StartServer(*listenPort, pro5)
}

func StartServer(port int, pro5 *Pro5Connection) {
  var mux = http.NewServeMux()
  mux.Handle("/connect", websocket.Handler(func(ws *websocket.Conn) { pro5.AddListener(ws) }))
  mux.Handle("/", http.FileServer(http.Dir("public/")))
  http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

func (c *Pro5Connection) SendUpdates(webSocket interface{}) {
//  conn <- c.DisplayLayouts
//  conn <- c.InitialSlide
//  c.OnSlide(func(newSlide interface{}) {
//    conn <- newSlide
//  })
}
//
//////////////////////////////

type Pro5Connection struct {
  Info Pro5ConnectInfo
  Connection net.Conn
  Listeners list.List
}

func ConnectToPro5(info Pro5ConnectInfo) (*Pro5Connection, error) {
  var result = new(Pro5Connection)
  var err error
  result.Info = info
  result.Connection, err = net.Dial("tcp", fmt.Sprintf("%s:%d", info.Host, info.Port))
  if err != nil {
    return nil, err
  }
  go result.Run()
  return result, nil
}

func (c *Pro5Connection) Run() {
  var xmlWriter = xml.NewEncoder(c.Connection)
  var loginElement = xml.StartElement{}
  loginElement.Name.Local = "StageDisplayLogin"
  xmlWriter.EncodeElement(c.Info.Password, loginElement)
  fmt.Fprintf(c.Connection, "\r\n")

  go c.ReadEverything()
}

func (c *Pro5Connection) ReadEverything() {
  var xmlReader = xml.NewDecoder(c.Connection)

  for {
    var token, err = xmlReader.Token()
    if err != nil {
      log.Fatal("xmlReader.Token(): ", err)
    }
    switch se := token.(type) {
    case xml.StartElement:
      fmt.Println(se.Name.Local)
    default:
      fmt.Printf("%T\n", token)
    }
  }
}

func (c *Pro5Connection) AddListener(listener io.Writer) {
  c.Listeners.PushBack(listener)
  fmt.Fprintf(listener, "Hello!")
}
