package main

import (
  "bytes"
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

type Pro5Connection struct {
  Info Pro5ConnectInfo
  Connection net.Conn
  Listeners list.List
  DisplayLayouts string
  LastSlide string
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

type DisplayLayout struct {
  // todo
}

type StageDisplayData struct {
  // todo
}

func (c *Pro5Connection) ReadEverything() {
  var xmlReader = xml.NewDecoder(c.Connection)

  // http://blog.davidsingleton.org/parsing-huge-xml-files-with-go/
  for {
    var token, err = xmlReader.Token()
    if err != nil {
      log.Fatal("xmlReader.Token(): ", err)
    }
    switch se := token.(type) {
    case xml.StartElement:
      switch se.Name.Local {
      case "DisplayLayouts":
        c.DisplayLayouts, err = readXmlString(xmlReader, DisplayLayout{}, &se)
        sendToListeners(c, c.DisplayLayouts)
      case "StageDisplayData":
        c.LastSlide, err = readXmlString(xmlReader, StageDisplayData{}, &se)
        sendToListeners(c, c.LastSlide)
      }
    }
  }
}

func readXmlString(xmlReader *xml.Decoder, v interface{}, startElement *xml.StartElement) (string, error) {
  var err = xmlReader.DecodeElement(v, startElement)
  if err != nil {
    return "", err
  }
  var buffer = new(bytes.Buffer)
  var xmlWriter = xml.NewEncoder(buffer)
  err = xmlWriter.Encode(v)
  if err != nil {
    return "", err
  }
  return buffer.String(), nil
}

func (c *Pro5Connection) AddListener(listener io.Writer) {
  c.Listeners.PushBack(listener)
  sendToListener(c, listener, c.DisplayLayouts)
  sendToListener(c, listener, c.LastSlide)
}

func sendToListeners(c *Pro5Connection, payload string) (err error) {
  for e := c.Listeners.Front(); e != nil; e = e.Next() {
    listener, ok := e.Value.(io.Writer)
    if ok {
      err = sendToListener(c, listener, payload)
      if err != nil {
        return
      }
    }
  }
  return
}

func sendToListener(c *Pro5Connection, listener io.Writer, payload string) error {
  fmt.Printf("%T %s\n", listener, payload)
  if len(payload) > 0 {
    _, err := fmt.Fprintf(listener, "%s", payload)
    if err != nil {
      return err
    }
  }
  return nil
}
